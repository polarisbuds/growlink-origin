package user_admin_usecase

import (
	"context"
	"errors"
	"fmt"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	"growlink/internal/shared/config"
	"growlink/internal/shared/crypto"
	logutils "growlink/internal/shared/logUtils"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"gorm.io/gorm"
)

type UserAdmin struct {
	userRepo       repository_interface.UserRepository
	deptMemberRepo repository_interface.DepartmentMemberRepository
	teamMemberRepo repository_interface.TeamMemberRepository
	deptRepo       repository_interface.DepartmentRepository
	teamRepo       repository_interface.TeamRepository
	roleRepo       repository_interface.RoleRepository
}

func NewUserAdmin(
	userRepo repository_interface.UserRepository,
	deptMemberRepo repository_interface.DepartmentMemberRepository,
	teamMemberRepo repository_interface.TeamMemberRepository,
	deptRepo repository_interface.DepartmentRepository,
	teamRepo repository_interface.TeamRepository,
	roleRepo repository_interface.RoleRepository,
) *UserAdmin {
	return &UserAdmin{
		userRepo:       userRepo,
		deptMemberRepo: deptMemberRepo,
		teamMemberRepo: teamMemberRepo,
		deptRepo:       deptRepo,
		teamRepo:       teamRepo,
		roleRepo:       roleRepo,
	}
}

func (uc *UserAdmin) ListUsers() (ret0 []outputmodel.ManagementUserOutput, err error) {
	const op = "user_admin_usecase.ListUsers"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	users, err := uc.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return uc.buildManagementUserOutputs(users)
}

type ManagementUsersListFilter struct {
	Name            string
	Email           string
	DepartmentUUIDs []uuid.UUID
	AccountStatus   string
}

func (f ManagementUsersListFilter) hasCondition() bool {
	if f.Name != "" {
		return true
	}
	if f.Email != "" {
		return true
	}
	if len(f.DepartmentUUIDs) > 0 {
		return true
	}
	if f.AccountStatus == "active" || f.AccountStatus == "suspended" {
		return true
	}
	return false
}

func (f ManagementUsersListFilter) suspendedFilter() *bool {
	switch f.AccountStatus {
	case "active":
		v := false
		return &v
	case "suspended":
		v := true
		return &v
	default:
		return nil
	}
}

func (uc *UserAdmin) ListUsersPaged(filter ManagementUsersListFilter, page, size int) (ret0 outputmodel.ManagementUsersPageOutput, err error) {
	const op = "user_admin_usecase.ListUsersPaged"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if page < 1 {
		page = 1
	}

	var users []entity.User
	if filter.hasCondition() {
		deptUUIDs := make([]string, len(filter.DepartmentUUIDs))
		for idx, departmentUUID := range filter.DepartmentUUIDs {
			deptUUIDs[idx] = departmentUUID.String()
		}
		hasDeptFilter := len(deptUUIDs) > 0
		deptMemberUserUUIDs := []string{}
		deptMemberUserUUIDSet := map[string]struct{}{}
		if hasDeptFilter {
			depts, err := uc.deptRepo.GetByUUIDs(deptUUIDs)
			if err != nil {
				return outputmodel.ManagementUsersPageOutput{}, err
			}
			deptIDs := make([]uint, len(depts))
			for idx, dept := range depts {
				deptIDs[idx] = dept.ID
			}

			departmentMemberByDeptUUID, err := uc.deptMemberRepo.GetByDepartmentIDs(deptIDs)
			if err != nil {
				return outputmodel.ManagementUsersPageOutput{}, err
			}

			for _, departmentMembers := range departmentMemberByDeptUUID {
				for _, member := range departmentMembers {
					if member.User == nil {
						continue
					}
					deptMemberUserUUIDSet[member.User.UUID.String()] = struct{}{}
				}
			}
			deptMemberUserUUIDs = make([]string, 0, len(deptMemberUserUUIDSet))
			for userUUID := range deptMemberUserUUIDSet {
				deptMemberUserUUIDs = append(deptMemberUserUUIDs, userUUID)
			}
		}
		users, err = uc.userRepo.GetByFilter(filter.Name, filter.Email, deptMemberUserUUIDs, filter.suspendedFilter())
	} else {
		users, err = uc.userRepo.GetAll()
	}
	if err != nil {
		return outputmodel.ManagementUsersPageOutput{}, err
	}

	pageUsers := sliceManagementUsersPage(users, page, size)
	items, err := uc.buildManagementUserOutputs(pageUsers)
	if err != nil {
		return outputmodel.ManagementUsersPageOutput{}, err
	}
	return outputmodel.ManagementUsersPageOutput{
		Items:    items,
		Total:    len(users),
		Page:     page,
		PageSize: size,
	}, nil
}

func (uc *UserAdmin) buildManagementUserOutputs(users []entity.User) ([]outputmodel.ManagementUserOutput, error) {
	group, _ := errgroup.WithContext(context.Background())

	roleIDToUUID := map[uint]string{}
	roleUUIDToName := map[string]string{}
	group.Go(func() error {
		roles, err := uc.roleRepo.GetAll()
		if err != nil {
			return err
		}
		for _, role := range roles {
			roleIDToUUID[role.ID] = role.UUID.String()
		}
		for _, role := range roles {
			roleUUIDToName[role.UUID.String()] = role.Name
		}
		return nil
	})

	deptIDToUUID := map[uint]string{}
	deptUUIDToName := map[string]string{}
	group.Go(func() error {
		depts, err := uc.deptRepo.GetAll()
		if err != nil {
			return err
		}
		for _, dept := range depts {
			u := dept.UUID.String()
			deptIDToUUID[dept.ID] = u
			deptUUIDToName[u] = dept.Name
		}
		return nil
	})

	out := make([]outputmodel.ManagementUserOutput, len(users))
	userIDs := make([]uint, len(users))
	for idx, user := range users {
		userIDs[idx] = user.ID
	}

	teamsMapByUserID, err := uc.loadTeamAssignmentsByUserIDs(userIDs)
	if err != nil {
		return nil, err
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	departmentsMapByUserID, err := uc.loadDepartmentAssignmentsByUserIDs(userIDs, deptIDToUUID, roleIDToUUID)
	if err != nil {
		return nil, err
	}

	sem := semaphore.NewWeighted(5)

	group, gctx := errgroup.WithContext(context.Background())
	ctx := context.WithoutCancel(gctx)
	for idx, u := range users {
		group.Go(func() error {
			if err := sem.Acquire(ctx, 1); err != nil {
				return err
			}
			defer sem.Release(1)

			departments := departmentsMapByUserID[u.ID]
			deptName, deptUUID, roleName, roleUUID := "", "", "", ""

			if len(departments) > 0 {
				deptUUID = departments[0].DepartmentUUID
				roleUUID = departments[0].RoleUUID
				deptName = deptUUIDToName[deptUUID]
				roleName = roleUUIDToName[roleUUID]
			}

			out[idx] = outputmodel.ManagementUserOutput{
				UUID:           u.UUID.String(),
				Name:           u.Name,
				Email:          u.Email,
				ImageURI:       u.ImageURI,
				Authority:      u.Authority.String(),
				SystemRole:     u.SystemRole.String(),
				DepartmentUUID: deptUUID,
				DepartmentName: deptName,
				RoleUUID:       roleUUID,
				RoleName:       roleName,
				Departments:    departments,
				Teams:          teamsMapByUserID[u.ID],
				Suspended:      u.Suspended,
			}
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}
	return out, nil
}

func (uc *UserAdmin) loadDepartmentAssignments(userID uint, deptIDToUUID map[uint]string, roleIDToUUID map[uint]string) ([]outputmodel.UserDepartmentAssignmentOutput, error) {
	members, err := uc.deptMemberRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	out := make([]outputmodel.UserDepartmentAssignmentOutput, 0, len(members))
	for _, member := range members {
		if member.DepartmentID == 0 || member.RoleID == 0 {
			continue
		}
		out = append(out, outputmodel.UserDepartmentAssignmentOutput{
			DepartmentName: member.Department.Name,
			DepartmentUUID: deptIDToUUID[member.DepartmentID],
			RoleUUID:       roleIDToUUID[member.RoleID],
		})
	}
	return out, nil
}

func (uc *UserAdmin) loadDepartmentAssignmentsByUserIDs(userIDs []uint, deptIDToUUID map[uint]string, roleIDToUUID map[uint]string) (map[uint][]outputmodel.UserDepartmentAssignmentOutput, error) {
	departmentMembersByUserID, err := uc.deptMemberRepo.GetByUserIDs(userIDs)
	if err != nil {
		return nil, err
	}

	departmentAssingmentByUserIDMap := map[uint][]outputmodel.UserDepartmentAssignmentOutput{}
	for userID, members := range departmentMembersByUserID {
		for _, member := range members {
			if member.DepartmentID == 0 || member.RoleID == 0 {
				continue
			}
			if _, ok := departmentAssingmentByUserIDMap[userID]; !ok {
				departmentAssingmentByUserIDMap[userID] = []outputmodel.UserDepartmentAssignmentOutput{}
			}
			departmentAssingmentByUserIDMap[userID] = append(departmentAssingmentByUserIDMap[userID], outputmodel.UserDepartmentAssignmentOutput{
				DepartmentName: member.Department.Name,
				DepartmentUUID: deptIDToUUID[member.DepartmentID],
				RoleUUID:       roleIDToUUID[member.RoleID],
			})
		}
	}
	return departmentAssingmentByUserIDMap, nil
}

func (uc *UserAdmin) loadTeamAssignments(userID uint) ([]outputmodel.UserTeamAssignmentOutput, error) {
	members, err := uc.teamMemberRepo.GetByMemberID(userID)
	if err != nil {
		return nil, err
	}
	teams, err := uc.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	teamIDToUUID := map[uint]string{}
	for _, team := range teams {
		teamIDToUUID[team.ID] = team.UUID.String()
	}
	out := make([]outputmodel.UserTeamAssignmentOutput, 0)
	seen := map[uint]struct{}{}
	for _, member := range members {
		if member.TeamID == 0 {
			continue
		}
		if _, ok := seen[member.TeamID]; ok {
			continue
		}
		seen[member.TeamID] = struct{}{}
		out = append(out, outputmodel.UserTeamAssignmentOutput{
			TeamUUID: teamIDToUUID[member.TeamID],
			IsLeader: member.IsReader,
		})
	}
	return out, nil
}

func (uc *UserAdmin) loadTeamAssignmentsByUserIDs(userIDs []uint) (map[uint][]outputmodel.UserTeamAssignmentOutput, error) {
	group, _ := errgroup.WithContext(context.Background())
	var err error

	teamMemberMapByMemberID := map[uint][]entity.TeamMember{}
	group.Go(func() error {
		teamMemberMapByMemberID, err = uc.teamMemberRepo.GetByMemberIDs(userIDs)
		if err != nil {
			return err
		}
		return nil
	})

	teams, err := uc.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	teamIDToUUID := map[uint]string{}
	for _, team := range teams {
		teamIDToUUID[team.ID] = team.UUID.String()
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	out := map[uint][]outputmodel.UserTeamAssignmentOutput{}
	for _, teamMembers := range teamMemberMapByMemberID {
		seen := map[uint]struct{}{}
		for _, teamMember := range teamMembers {
			if teamMember.TeamID == 0 {
				continue
			}
			if _, ok := seen[teamMember.TeamID]; ok {
				continue
			}
			seen[teamMember.TeamID] = struct{}{}
			if _, ok := out[teamMember.MemberID]; !ok {
				out[teamMember.MemberID] = []outputmodel.UserTeamAssignmentOutput{}
			}

			out[teamMember.MemberID] = append(out[teamMember.MemberID], outputmodel.UserTeamAssignmentOutput{
				TeamUUID: teamIDToUUID[teamMember.TeamID],
				IsLeader: teamMember.IsReader,
			})
		}
	}
	return out, nil
}

type DepartmentAssignmentInput struct {
	DepartmentUuid string `json:"departmentUuid"`
	RoleUuid       string `json:"roleUuid"`
}

type TeamAssignmentInput struct {
	TeamUuid string `json:"teamUuid"`
}

type CreateUserInput struct {
	Name        string                      `json:"name"`
	Email       string                      `json:"email"`
	Departments []DepartmentAssignmentInput `json:"departments"`
	Teams       []TeamAssignmentInput       `json:"teams"`
	Password    string                      `json:"password"`
}

func (uc *UserAdmin) CreateUser(input CreateUserInput) (ret0 *outputmodel.ManagementUserOutput, err error) {
	const op = "user_admin_usecase.CreateUser"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(input.Email)
	if name == "" {
		return nil, fmt.Errorf("氏名を入力してください")
	}
	if email == "" {
		return nil, fmt.Errorf("メールを入力してください")
	}
	if len(input.Departments) == 0 {
		return nil, fmt.Errorf("所属部署を選択してください")
	}
	departments, err := uc.resolveDepartmentAssignments(input.Departments)
	if err != nil {
		return nil, err
	}
	teamsResolved, err := uc.resolveTeamAssignments(input.Teams)
	if err != nil {
		return nil, err
	}
	existing, err := uc.userRepo.GetByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("このメールアドレスは既に登録されています")
	}
	user := &entity.User{
		Name:       name,
		Email:      email,
		Authority:  valueobject.AuthorityGeneral,
		SystemRole: valueobject.SystemRoleNone,
		Suspended:  false,
	}
	if strings.TrimSpace(input.Password) != "" {
		hashed, err := crypto.HashPassword(strings.TrimSpace(input.Password), config.GetPepper())
		if err != nil {
			return nil, err
		}
		user.Password = &hashed
	}
	if err = uc.userRepo.Create(user); err != nil {
		return nil, err
	}
	if err = uc.syncAffiliations(user.ID, departments, teamsResolved); err != nil {
		return nil, err
	}
	teams, err := uc.loadTeamAssignments(user.ID)
	if err != nil {
		return nil, err
	}
	listed, err := uc.departmentAssignmentsOutput(departments)
	if err != nil {
		return nil, err
	}
	deptUUID, roleUUID := listed[0].DepartmentUUID, listed[0].RoleUUID
	return &outputmodel.ManagementUserOutput{
		UUID:           user.UUID.String(),
		Name:           user.Name,
		Email:          user.Email,
		Authority:      user.Authority.String(),
		SystemRole:     user.SystemRole.String(),
		DepartmentUUID: deptUUID,
		RoleUUID:       roleUUID,
		Departments:    listed,
		Teams:          teams,
		Suspended:      user.Suspended,
	}, nil
}

type UpdateUserInput struct {
	Authority   string                      `json:"authority"`
	SystemRole  string                      `json:"systemRole"`
	Departments []DepartmentAssignmentInput `json:"departments"`
	Teams       []TeamAssignmentInput       `json:"teams"`
	Suspended   *bool                       `json:"suspended"`
}

func (uc *UserAdmin) UpdateUser(userID uint, input UpdateUserInput) (err error) {
	const op = "user_admin_usecase.UpdateUser"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if len(input.Departments) == 0 {
		return fmt.Errorf("所属部署を選択してください")
	}
	departments, err := uc.resolveDepartmentAssignments(input.Departments)
	if err != nil {
		return err
	}
	teamsResolved, err := uc.resolveTeamAssignments(input.Teams)
	if err != nil {
		return err
	}
	columns, err := userColumnChanges(user, input, departments)
	if err != nil {
		return err
	}
	if len(columns) > 0 {
		if err = uc.userRepo.UpdateColumns(user.ID, columns); err != nil {
			return err
		}
	}
	return uc.syncAffiliations(userID, departments, teamsResolved)
}

func userColumnChanges(user *entity.User, input UpdateUserInput, departments []deptAssignmentResolved) (map[string]any, error) {
	columns := map[string]any{}
	if input.Authority != "" {
		auth, err := valueobject.ParseAuthority(input.Authority)
		if err != nil {
			return nil, err
		}
		if user.Authority != auth {
			columns["authority"] = auth.Int()
		}
	}
	if input.SystemRole != "" {
		role, err := valueobject.ParseSystemRole(input.SystemRole)
		if err != nil {
			return nil, err
		}
		if user.SystemRole != role {
			columns["system_role"] = role.Int()
		}
	}
	if input.Suspended != nil && user.Suspended != *input.Suspended {
		columns["suspended"] = *input.Suspended
	}
	return columns, nil
}

func diffIDs(existing, desired []uint) (remove, add []uint) {
	have := map[uint]struct{}{}
	for _, id := range existing {
		have[id] = struct{}{}
	}
	want := map[uint]struct{}{}
	for _, id := range desired {
		want[id] = struct{}{}
	}
	for id := range have {
		if _, ok := want[id]; !ok {
			remove = append(remove, id)
		}
	}
	for id := range want {
		if _, ok := have[id]; !ok {
			add = append(add, id)
		}
	}
	return remove, add
}

func (uc *UserAdmin) syncAffiliations(userID uint, departments []deptAssignmentResolved, teams []teamAssignmentResolved) error {
	allowed, err := uc.syncDepartmentAssignments(userID, departments)
	if err != nil {
		return err
	}
	return uc.syncTeamMemberships(userID, teams, allowed)
}

func (uc *UserAdmin) syncDepartmentAssignments(userID uint, rows []deptAssignmentResolved) (map[uint]struct{}, error) {
	roles, err := uc.roleRepo.GetAll()
	if err != nil {
		return nil, err
	}
	roleNameByID := map[uint]string{}
	for _, role := range roles {
		roleNameByID[role.ID] = role.Name
	}
	want := make([]struct {
		departmentID uint
		roleID       uint
	}, 0, len(rows))
	allowed := map[uint]struct{}{}
	for _, row := range rows {
		if row.DepartmentID == 0 {
			return nil, fmt.Errorf("所属部署を選択してください")
		}
		if _, dup := allowed[row.DepartmentID]; dup {
			return nil, fmt.Errorf("同じ部署を重複して登録できません")
		}
		if _, ok := roleNameByID[row.RoleID]; !ok {
			return nil, fmt.Errorf("役職を選択してください")
		}
		allowed[row.DepartmentID] = struct{}{}
		want = append(want, struct {
			departmentID uint
			roleID       uint
		}{
			departmentID: row.DepartmentID,
			roleID:       row.RoleID,
		})
	}

	membersOfUser, err := uc.deptMemberRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	previous := map[uint]struct{}{}
	for _, member := range membersOfUser {
		previous[member.DepartmentID] = struct{}{}
	}

	syncRows := make([]entity.DepartmentMember, len(want))
	for i, row := range want {
		syncRows[i] = entity.DepartmentMember{
			DepartmentID: row.departmentID,
			UserID:       userID,
			RoleID:       row.roleID,
		}
	}
	if err := uc.deptMemberRepo.SyncUserDepartments(userID, syncRows); err != nil {
		return nil, err
	}

	removedIDs := make([]uint, 0)
	for id := range previous {
		if _, keep := allowed[id]; !keep {
			removedIDs = append(removedIDs, id)
		}
	}
	if len(removedIDs) == 0 {
		return allowed, nil
	}
	depts, err := uc.deptRepo.GetByIDs(removedIDs)
	if err != nil {
		return nil, err
	}
	for _, dept := range depts {
		if dept.HeadUserID == userID {
			if err := uc.deptRepo.UpdateHeadUserID(dept.ID, 0); err != nil {
				return nil, err
			}
		}
	}

	return allowed, nil
}

func (uc *UserAdmin) syncTeamMemberships(userID uint, resolvedTeams []teamAssignmentResolved, allowedDepartments map[uint]struct{}) error {
	want := map[uint]struct{}{}
	teamIDs := []uint{}
	for _, row := range resolvedTeams {
		teamIDs = append(teamIDs, row.TeamID)
	}

	teams, err := uc.teamRepo.GetByIDs(teamIDs)
	if err != nil {
		return fmt.Errorf("チームの取得に失敗しました: %w", err)
	}

	for _, team := range teams {
		if _, ok := allowedDepartments[team.DepartmentID]; !ok {
			return fmt.Errorf("選択した部署に属さないチームは登録できません")
		}
		want[team.ID] = struct{}{}
	}

	existing, err := uc.teamMemberRepo.GetByMemberID(userID)
	if err != nil {
		return err
	}
	have := map[uint]struct{}{}
	for _, member := range existing {
		if _, ok := want[member.TeamID]; !ok {
			if err := uc.teamMemberRepo.Delete(member.ID); err != nil {
				return err
			}
			continue
		}
		have[member.TeamID] = struct{}{}
	}
	for teamID := range want {
		if _, ok := have[teamID]; ok {
			continue
		}
		if err := uc.teamMemberRepo.Create(teamID, userID, false); err != nil {
			return err
		}
	}
	return nil
}

func (uc *UserAdmin) SetUsersSuspended(uuids []string, suspended bool) (err error) {
	const op = "user_admin_usecase.SetUsersSuspended"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if len(uuids) == 0 {
		return fmt.Errorf("対象のユーザが選択されていません")
	}

	users, err := uc.userRepo.GetByUUIDs(uuids)
	if err != nil {
		return err
	}
	ids := make([]uint, len(users))
	for idx, user := range users {
		ids[idx] = user.ID
	}
	return uc.userRepo.SetSuspendedByIDs(ids, suspended)
}

func ValidateLineWorksLocalID(id string) error {
	if strings.Contains(id, "@") {
		return fmt.Errorf("LINE WORKS ID に @ を含めることはできません")
	}
	return nil
}
