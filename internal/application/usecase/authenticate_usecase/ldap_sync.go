package authenticate_usecase

import (
	"context"
	"log/slog"
	"sort"
	"strings"
	"time"

	"growlink/internal/domain/entity"
	ldap_entity "growlink/internal/domain/entity/ldap"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"
)

func (uc *AuthUsecase) SyncDirectory(ctx context.Context) (err error) {
	const op = opAuth + ".SyncDirectory"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if uc.ldapClient == nil {
		logutils.Warn(ctx, logutils.LayerUsecase, op, "LDAPクライアントが無いため同期をスキップ")
		return nil
	}

	listed, err := uc.ldapClient.ListProperActive()
	if err != nil {
		logutils.Warn(ctx, logutils.LayerUsecase, op, "LDAP一覧の取得に失敗したため同期を中断", slog.String("err", err.Error()))
		return nil
	}
	if len(listed) == 0 {
		logutils.Warn(ctx, logutils.LayerUsecase, op, "LDAP一覧が0件のため同期を中断")
		return nil
	}

	users, err := uc.userRepo.GetAll()
	if err != nil {
		return err
	}
	byEmail := map[string]*entity.User{}
	for i := range users {
		key := normalizeEmail(users[i].Email)
		if key == "" {
			continue
		}
		if _, ok := byEmail[key]; !ok {
			byEmail[key] = &users[i]
		}
	}

	seen := map[string]struct{}{}
	for _, ldapUser := range listed {
		key := normalizeEmail(ldapUser.Mail)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		user, ok := byEmail[key]
		if !ok {
			created, err := uc.createUserFromLDAP(ctx, ldapUser)
			if err != nil {
				return err
			}
			byEmail[key] = created
			continue
		}
		if userHasPassword(user) {
			continue
		}
		if err = uc.refreshLDAPUser(ctx, user, ldapUser); err != nil {
			return err
		}
	}

	var suspendIDs []uint
	for i := range users {
		user := &users[i]
		if user.Suspended || userHasPassword(user) {
			continue
		}
		if _, ok := seen[normalizeEmail(user.Email)]; ok {
			continue
		}
		suspendIDs = append(suspendIDs, user.ID)
	}
	if len(suspendIDs) == 0 {
		return nil
	}
	if err = uc.userRepo.SetSuspendedByIDs(suspendIDs, true); err != nil {
		return err
	}
	logutils.Info(ctx, logutils.LayerUsecase, op, "LDAPに存在しないユーザを停止", slog.Int("count", len(suspendIDs)))
	return nil
}

func (uc *AuthUsecase) createUserFromLDAP(ctx context.Context, ldapUser ldap_entity.LDAPUser) (*entity.User, error) {
	if _, err := uc.ensureDepartmentFromLDAP(ctx, ldapUser.Department); err != nil {
		return nil, err
	}

	var lineworksID string
	email := strings.TrimSpace(ldapUser.Mail)
	atMarkIdx := strings.Index(ldapUser.Mail, "@")
	if atMarkIdx >= 0 {
		lineworksID = ldapUser.Mail[:atMarkIdx] + "@cap"
	}
	user := &entity.User{
		Name:        strings.TrimSpace(ldapUser.Name),
		Email:       email,
		Authority:   valueobject.AuthorityGeneral,
		LineWorksID: lineworksID,
	}
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}
	deptID, err := uc.departmentIDForNewLDAPUser(ldapUser.Department)
	if err != nil {
		return nil, err
	}
	member := entity.DepartmentMember{
		DepartmentID: deptID,
		UserID:       user.ID,
		RoleID:       defaultRoleID(),
	}
	if err := uc.deptMemberRepo.ReplaceByUser(user.ID, member); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *AuthUsecase) refreshLDAPUser(ctx context.Context, user *entity.User, ldapUser ldap_entity.LDAPUser) error {
	name := strings.TrimSpace(ldapUser.Name)
	if name != "" && user.Name != name {
		if err := uc.userRepo.UpdateColumns(user.ID, map[string]any{"name": name}); err != nil {
			return err
		}
		user.Name = name
	}
	if strings.TrimSpace(ldapUser.Department) == "" {
		return nil
	}
	deptID, err := uc.ensureDepartmentFromLDAP(ctx, ldapUser.Department)
	if err != nil {
		return err
	}
	if deptID == 0 {
		return nil
	}
	desired := membershipForLDAPDepartment(user.ID, deptID, user.DepartmentMembers)
	if sameSingleMembership(user.DepartmentMembers, desired) {
		return nil
	}
	return uc.deptMemberRepo.ReplaceByUser(user.ID, desired)
}

func membershipForLDAPDepartment(userID, deptID uint, members []entity.DepartmentMember) entity.DepartmentMember {
	ordered := append([]entity.DepartmentMember(nil), members...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].ID < ordered[j].ID })

	desired := entity.DepartmentMember{
		DepartmentID: deptID,
		UserID:       userID,
		RoleID:       defaultRoleID(),
	}
	for _, current := range ordered {
		if current.DepartmentID != deptID {
			continue
		}
		desired.RoleID = current.RoleID
		desired.IsHead = current.IsHead
		desired.IsManager = current.IsManager
		return desired
	}
	if len(ordered) > 0 && ordered[0].RoleID > 0 {
		desired.RoleID = ordered[0].RoleID
	}
	return desired
}

func sameSingleMembership(members []entity.DepartmentMember, desired entity.DepartmentMember) bool {
	if len(members) != 1 {
		return false
	}
	current := members[0]
	return current.DepartmentID == desired.DepartmentID &&
		current.RoleID == desired.RoleID &&
		current.IsHead == desired.IsHead &&
		current.IsManager == desired.IsManager
}

func userHasPassword(user *entity.User) bool {
	return user != nil && user.Password != nil && strings.TrimSpace(*user.Password) != ""
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func NextDirectorySyncAt(now time.Time) time.Time {
	loc := now.Location()
	next := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, loc)
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}
