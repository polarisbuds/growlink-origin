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
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
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
	users, err := uc.userRepo.GetAll()
	if err != nil {
		return outputmodel.ManagementUsersPageOutput{}, err
	}
	deptUUIDSet := map[string]struct{}{}

	for _, d := range filter.DepartmentUUIDs {
		deptUUIDSet[d.String()] = struct{}{}
	}
	hasDeptFilter := len(deptUUIDSet) > 0
	var deptMemberUserIDs map[uint]struct{}
	if hasDeptFilter {
		depts, err := uc.deptRepo.GetAll()
		if err != nil {
			return outputmodel.ManagementUsersPageOutput{}, err
		}
		deptIDToUUID := map[uint]string{}
		for _, dept := range depts {
			deptIDToUUID[dept.ID] = dept.UUID.String()
		}
		members, err := uc.deptMemberRepo.GetAll()
		if err != nil {
			return outputmodel.ManagementUsersPageOutput{}, err
		}
		deptMemberUserIDs = userIDsInDepartments(members, deptUUIDSet, deptIDToUUID)
	}
	filtered := filterManagementUsers(users, filter.Name, filter.Email, deptMemberUserIDs, hasDeptFilter)
	pageUsers := sliceManagementUsersPage(filtered, page, size)
	items, err := uc.buildManagementUserOutputs(pageUsers)
	if err != nil {
		return outputmodel.ManagementUsersPageOutput{}, err
	}
	return outputmodel.ManagementUsersPageOutput{
		Items:    items,
		Total:    len(filtered),
		Page:     page,
		PageSize: size,
	}, nil
}

func (uc *UserAdmin) buildManagementUserOutputs(users []entity.User) ([]outputmodel.ManagementUserOutput, error) {
	roles, err := uc.roleRepo.GetAll()
	if err != nil {
		return nil, err
	}
	roleIDToUUID := make(map[uint]string, len(roles))
	roleUUIDToName := make(map[string]string, len(roles))
	roleWg := sync.WaitGroup{}
	roleWg.Go(func() {
		for _, role := range roles {
			roleIDToUUID[role.ID] = role.UUID.String()
		}

		for _, role := range roles {
			roleUUIDToName[role.UUID.String()] = role.Name
		}
	})

	depts, err := uc.deptRepo.GetAll()
	if err != nil {
		return nil, err
	}
	deptIDToUUID := make(map[uint]string, len(depts))
	deptUUIDToName := make(map[string]string, len(depts))
	deptWg := sync.WaitGroup{}
	deptWg.Go(func() {
		for _, dept := range depts {
			u := dept.UUID.String()
			deptIDToUUID[dept.ID] = u
			deptUUIDToName[u] = dept.Name
		}
	})
	roleWg.Wait()
	deptWg.Wait()

	group, _ := errgroup.WithContext(context.Background())
	out := make([]outputmodel.ManagementUserOutput, len(users))
	for idx, u := range users {
		group.Go(func() error {
			departments, err := uc.loadDepartmentAssignments(u.ID, deptIDToUUID, roleIDToUUID)
			if err != nil {
				return err
			}
			teams, err := uc.loadTeamAssignments(u.ID)
			if err != nil {
				return err
			}
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
				Teams:          teams,
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

	previous := map[uint]struct{}{}
	membersOfUser, err := uc.deptMemberRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	for _, member := range membersOfUser {
		previous[member.DepartmentID] = struct{}{}
	}
	touched := map[uint]struct{}{}
	for id := range allowed {
		touched[id] = struct{}{}
	}
	for id := range previous {
		touched[id] = struct{}{}
	}
	roleIDByDept := map[uint]uint{}
	for _, row := range want {
		roleIDByDept[row.departmentID] = row.roleID
	}

	// TODO; members, err := uc.deptMemberRepo.GetByDepartmentID(deptID)これをループ外へ出したい

	departments, err := uc.deptRepo.GetByIDsByUnscoped(slices.Collect(maps.Keys(touched)))
	if err != nil {
		return nil, err
	}

	for _, dept := range departments {
		_, keep := allowed[dept.ID]
		next := make([]entity.DepartmentMember, 0, len(dept.DepartmentMembers)+1)
		found := false
		removed := false
		for _, member := range dept.DepartmentMembers {
			if member.UserID == userID {
				if !keep {
					removed = true
					continue
				}
				found = true
				if roleID, ok := roleIDByDept[dept.ID]; ok {
					member.RoleID = roleID
				}
				next = append(next, member)
				continue
			}
			next = append(next, member)
		}
		if keep && !found {
			roleID := roleIDByDept[dept.ID]
			if roleID == 0 {
				return nil, fmt.Errorf("役職を選択してください")
			}
			next = append(next, entity.DepartmentMember{
				DepartmentID: dept.ID,
				UserID:       userID,
				RoleID:       roleID,
			})
		}
		if !departmentMembersEqual(dept.DepartmentMembers, next) {
			if err := uc.deptMemberRepo.ReplaceByDepartment(dept.ID, next); err != nil {
				return nil, err
			}
		}
		if removed {
			dept, err := uc.deptRepo.GetByID(dept.ID)
			if err != nil {
				return nil, err
			}
			if dept.HeadUserID == userID {
				if err := uc.deptRepo.UpdateHeadUserID(dept.ID, 0); err != nil {
					return nil, err
				}
			}
		}
	}

	return allowed, nil
}

func departmentMembersEqual(current, next []entity.DepartmentMember) bool {
	count := map[[4]uint]int{}
	n := 0
	for _, member := range current {
		if member.UserID == 0 {
			continue
		}
		n++
		count[memberKey(member)]++
	}
	for _, member := range next {
		if member.UserID == 0 {
			continue
		}
		key := memberKey(member)
		if count[key] == 0 {
			return false
		}
		count[key]--
		n--
	}
	return n == 0
}

func memberKey(member entity.DepartmentMember) [4]uint {
	head := uint(0)
	manager := uint(0)
	if member.IsHead {
		head = 1
	}
	if member.IsManager {
		manager = 1
	}
	return [4]uint{member.UserID, head, manager, member.RoleID}
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

	ids := make([]uint, 0, len(uuids))
	users, err := uc.userRepo.GetByUUIDs(uuids)
	if err != nil {
		return err
	}
	for _, user := range users {
		ids = append(ids, user.ID)
	}
	return uc.userRepo.SetSuspendedByIDs(ids, suspended)
}

func ValidateLineWorksLocalID(id string) error {
	if strings.Contains(id, "@") {
		return fmt.Errorf("LINE WORKS ID に @ を含めることはできません")
	}
	return nil
}
