package user_admin_usecase

import (
	"errors"
	"fmt"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	"growlink/internal/shared/config"
	"growlink/internal/shared/crypto"
	"strings"

	"gorm.io/gorm"
)

type UserAdmin struct {
	userRepo         repository_interface.UserRepository
	deptRefRepo      repository_interface.UserReferenceDepartmentRepository
	teamRefRepo      repository_interface.UserReferenceTeamRepository
	relationshipRepo repository_interface.UserRelationshipRepository
	userDeptRepo     repository_interface.UserDepartmentRepository
	deptRoleRepo     repository_interface.DepartmentRoleRepository
	deptMemberRepo   repository_interface.DepartmentMemberRepository
	teamMemberRepo   repository_interface.TeamMemberRepository
	deptRepo         repository_interface.DepartmentRepository
	teamRepo         repository_interface.TeamRepository
	roleRepo         repository_interface.RoleRepository
}

func NewUserAdmin(
	userRepo repository_interface.UserRepository,
	deptRefRepo repository_interface.UserReferenceDepartmentRepository,
	teamRefRepo repository_interface.UserReferenceTeamRepository,
	relationshipRepo repository_interface.UserRelationshipRepository,
	userDeptRepo repository_interface.UserDepartmentRepository,
	deptRoleRepo repository_interface.DepartmentRoleRepository,
	deptMemberRepo repository_interface.DepartmentMemberRepository,
	teamMemberRepo repository_interface.TeamMemberRepository,
	deptRepo repository_interface.DepartmentRepository,
	teamRepo repository_interface.TeamRepository,
	roleRepo repository_interface.RoleRepository,
) *UserAdmin {
	return &UserAdmin{
		userRepo:         userRepo,
		deptRefRepo:      deptRefRepo,
		teamRefRepo:      teamRefRepo,
		relationshipRepo: relationshipRepo,
		userDeptRepo:     userDeptRepo,
		deptRoleRepo:     deptRoleRepo,
		deptMemberRepo:   deptMemberRepo,
		teamMemberRepo:   teamMemberRepo,
		deptRepo:         deptRepo,
		teamRepo:         teamRepo,
		roleRepo:         roleRepo,
	}
}

func (uc *UserAdmin) ListUsers() ([]outputmodel.ManagementUserOutput, error) {
	users, err := uc.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	roles, err := uc.roleRepo.GetAll()
	if err != nil {
		return nil, err
	}
	roleIDByName := map[string]uint{}
	for _, role := range roles {
		roleIDByName[role.Name] = role.ID
	}
	out := make([]outputmodel.ManagementUserOutput, 0, len(users))
	for _, u := range users {
		targetIDs, err := uc.loadTargetIDs(u.ID, u.OperationTarget)
		if err != nil {
			return nil, err
		}
		departments, err := uc.loadDepartmentAssignments(u, roleIDByName)
		if err != nil {
			return nil, err
		}
		teams, err := uc.loadTeamAssignments(u.ID)
		if err != nil {
			return nil, err
		}
		deptName := ""
		if u.Department != nil {
			deptName = u.Department.Name
		}
		roleName := ""
		if u.Role != nil {
			roleName = u.Role.Name
		}
		out = append(out, outputmodel.ManagementUserOutput{
			ID:              u.ID,
			Name:            u.Name,
			Email:           u.Email,
			ImageURI:        u.ImageURI,
			Authority:       u.Authority.String(),
			SystemRole:      u.SystemRole.String(),
			OperationTarget: u.OperationTarget.String(),
			TargetIDs:       targetIDs,
			DepartmentID:    u.DepartmentID,
			DepartmentName:  deptName,
			RoleID:          u.RoleID,
			RoleName:        roleName,
			Departments:     departments,
			Teams:           teams,
			Suspended:       u.Suspended,
		})
	}
	return out, nil
}

func (uc *UserAdmin) loadDepartmentAssignments(user entity.User, roleIDByName map[string]uint) ([]outputmodel.UserDepartmentAssignmentOutput, error) {
	roles, err := uc.deptRoleRepo.GetByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	out := make([]outputmodel.UserDepartmentAssignmentOutput, 0, len(roles))
	seen := map[uint]struct{}{}
	for _, role := range roles {
		if role.DepartmentID == 0 {
			continue
		}
		if _, ok := seen[role.DepartmentID]; ok {
			continue
		}
		seen[role.DepartmentID] = struct{}{}
		out = append(out, outputmodel.UserDepartmentAssignmentOutput{
			DepartmentID: role.DepartmentID,
			RoleID:       roleIDByName[role.Name],
		})
	}
	if len(out) == 0 && user.DepartmentID != 0 {
		out = append(out, outputmodel.UserDepartmentAssignmentOutput{
			DepartmentID: user.DepartmentID,
			RoleID:       user.RoleID,
		})
	}
	return out, nil
}

func (uc *UserAdmin) loadTeamAssignments(userID uint) ([]outputmodel.UserTeamAssignmentOutput, error) {
	members, err := uc.teamMemberRepo.GetByMemberID(userID)
	if err != nil {
		return nil, err
	}
	out := make([]outputmodel.UserTeamAssignmentOutput, 0)
	seen := map[uint]struct{}{}
	for _, member := range members {
		if !member.IsReader || member.TeamID == 0 {
			continue
		}
		if _, ok := seen[member.TeamID]; ok {
			continue
		}
		seen[member.TeamID] = struct{}{}
		out = append(out, outputmodel.UserTeamAssignmentOutput{TeamID: member.TeamID})
	}
	return out, nil
}

func (uc *UserAdmin) loadTargetIDs(userID uint, target valueobject.OperationTarget) ([]uint, error) {
	switch target {
	case valueobject.OperationTargetDepartment:
		refs, err := uc.deptRefRepo.GetByUserID(userID)
		if err != nil {
			return nil, err
		}
		ids := make([]uint, 0, len(refs))
		for _, r := range refs {
			ids = append(ids, r.ReferableDepartmentID)
		}
		return ids, nil
	case valueobject.OperationTargetTeam:
		refs, err := uc.teamRefRepo.GetByUserID(userID)
		if err != nil {
			return nil, err
		}
		ids := make([]uint, 0, len(refs))
		for _, r := range refs {
			ids = append(ids, r.ReferableTeamID)
		}
		return ids, nil
	case valueobject.OperationTargetSubordinate:
		rels, err := uc.relationshipRepo.GetByBossID(userID)
		if err != nil {
			return nil, err
		}
		ids := make([]uint, 0, len(rels))
		for _, r := range rels {
			ids = append(ids, r.SourceUserID)
		}
		return ids, nil
	default:
		return []uint{}, nil
	}
}

type DepartmentAssignmentInput struct {
	DepartmentID uint `json:"departmentId"`
	RoleID       uint `json:"roleId"`
}

type TeamAssignmentInput struct {
	TeamID uint `json:"teamId"`
}

type CreateUserInput struct {
	Name         string                      `json:"name"`
	Email        string                      `json:"email"`
	DepartmentID uint                        `json:"departmentId"`
	RoleID       uint                        `json:"roleId"`
	Departments  []DepartmentAssignmentInput `json:"departments"`
	Teams        []TeamAssignmentInput       `json:"teams"`
	Password     string                      `json:"password"`
}

func (uc *UserAdmin) CreateUser(input CreateUserInput) (*outputmodel.ManagementUserOutput, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(input.Email)
	if name == "" {
		return nil, fmt.Errorf("氏名を入力してください")
	}
	if email == "" {
		return nil, fmt.Errorf("メールを入力してください")
	}
	departments := input.Departments
	if len(departments) == 0 && input.DepartmentID != 0 && input.RoleID != 0 {
		departments = []DepartmentAssignmentInput{{
			DepartmentID: input.DepartmentID,
			RoleID:       input.RoleID,
		}}
	}
	if len(departments) == 0 {
		return nil, fmt.Errorf("所属部署を選択してください")
	}
	existing, err := uc.userRepo.GetByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("このメールアドレスは既に登録されています")
	}
	user := &entity.User{
		Name:            name,
		Email:           email,
		DepartmentID:    departments[0].DepartmentID,
		RoleID:          departments[0].RoleID,
		Authority:       valueobject.AuthorityGeneral,
		SystemRole:      valueobject.SystemRoleNone,
		OperationTarget: valueobject.OperationTargetNone,
		Suspended:       false,
	}
	if strings.TrimSpace(input.Password) != "" {
		hashed, err := crypto.HashPassword(strings.TrimSpace(input.Password), config.GetPepper())
		if err != nil {
			return nil, err
		}
		user.Password = &hashed
	}
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}
	if err := uc.syncAffiliations(user.ID, departments, input.Teams); err != nil {
		return nil, err
	}
	teams, err := uc.loadTeamAssignments(user.ID)
	if err != nil {
		return nil, err
	}
	listed := make([]outputmodel.UserDepartmentAssignmentOutput, 0, len(departments))
	for _, row := range departments {
		listed = append(listed, outputmodel.UserDepartmentAssignmentOutput{
			DepartmentID: row.DepartmentID,
			RoleID:       row.RoleID,
		})
	}
	return &outputmodel.ManagementUserOutput{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		Authority:       user.Authority.String(),
		SystemRole:      user.SystemRole.String(),
		OperationTarget: user.OperationTarget.String(),
		TargetIDs:       []uint{},
		DepartmentID:    user.DepartmentID,
		RoleID:          user.RoleID,
		Departments:     listed,
		Teams:           teams,
		Suspended:       user.Suspended,
	}, nil
}

type UpdateUserInput struct {
	Authority       string                      `json:"authority"`
	SystemRole      string                      `json:"systemRole"`
	OperationTarget string                      `json:"operationTarget"`
	TargetIDs       []uint                      `json:"targetIds"`
	DepartmentID    uint                        `json:"departmentId"`
	RoleID          uint                        `json:"roleId"`
	Departments     []DepartmentAssignmentInput `json:"departments"`
	Teams           []TeamAssignmentInput       `json:"teams"`
	Suspended       *bool                       `json:"suspended"`
}

func (uc *UserAdmin) UpdateUser(userID uint, input UpdateUserInput) error {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if input.Authority != "" {
		auth, err := valueobject.ParseAuthority(input.Authority)
		if err != nil {
			return err
		}
		user.Authority = auth
	}
	if input.SystemRole != "" {
		role, err := valueobject.ParseSystemRole(input.SystemRole)
		if err != nil {
			return err
		}
		user.SystemRole = role
	}
	if input.OperationTarget != "" {
		target, err := valueobject.ParseOperationTarget(input.OperationTarget)
		if err != nil {
			return err
		}
		user.OperationTarget = target
	}
	departments := input.Departments
	if len(departments) == 0 && input.DepartmentID != 0 && input.RoleID != 0 {
		departments = []DepartmentAssignmentInput{{
			DepartmentID: input.DepartmentID,
			RoleID:       input.RoleID,
		}}
	}
	if len(departments) == 0 {
		return fmt.Errorf("所属部署を選択してください")
	}
	user.DepartmentID = departments[0].DepartmentID
	user.RoleID = departments[0].RoleID
	if input.Suspended != nil {
		user.Suspended = *input.Suspended
	}
	if err := uc.userRepo.UpdateByID(user); err != nil {
		return err
	}
	if err := uc.syncAffiliations(userID, departments, input.Teams); err != nil {
		return err
	}
	return uc.syncOperationTargets(userID, user.OperationTarget, input.TargetIDs)
}

func (uc *UserAdmin) syncAffiliations(userID uint, departments []DepartmentAssignmentInput, teams []TeamAssignmentInput) error {
	allowed, err := uc.syncDepartmentAssignments(userID, departments)
	if err != nil {
		return err
	}
	return uc.syncTeamLeaders(userID, teams, allowed)
}

func (uc *UserAdmin) syncDepartmentAssignments(userID uint, rows []DepartmentAssignmentInput) (map[uint]struct{}, error) {
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
		roleName     string
		isHead       bool
	}, 0, len(rows))
	allowed := map[uint]struct{}{}
	for _, row := range rows {
		if row.DepartmentID == 0 {
			return nil, fmt.Errorf("所属部署を選択してください")
		}
		if _, dup := allowed[row.DepartmentID]; dup {
			return nil, fmt.Errorf("同じ部署を重複して登録できません")
		}
		roleName, ok := roleNameByID[row.RoleID]
		if !ok || row.RoleID == 0 {
			return nil, fmt.Errorf("役職を選択してください")
		}
		allowed[row.DepartmentID] = struct{}{}
		want = append(want, struct {
			departmentID uint
			roleName     string
			isHead       bool
		}{
			departmentID: row.DepartmentID,
			roleName:     roleName,
			isHead:       roleName == "部長",
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
	deptIDs, err := uc.userDeptRepo.GetDepartmentIDsByUserID(userID)
	if err != nil {
		return nil, err
	}
	for _, id := range deptIDs {
		previous[id] = struct{}{}
	}

	touched := map[uint]struct{}{}
	for id := range allowed {
		touched[id] = struct{}{}
	}
	for id := range previous {
		touched[id] = struct{}{}
	}
	headByDept := map[uint]bool{}
	for _, row := range want {
		headByDept[row.departmentID] = row.isHead
	}

	for deptID := range touched {
		members, err := uc.deptMemberRepo.GetByDepartmentID(deptID)
		if err != nil {
			return nil, err
		}
		dept, err := uc.deptRepo.GetByID(deptID)
		if err != nil {
			return nil, err
		}
		_, keep := allowed[deptID]
		isHead := headByDept[deptID]
		next := make([]entity.DepartmentMember, 0, len(members)+1)
		found := false
		for _, member := range members {
			if member.UserID == userID {
				if !keep {
					continue
				}
				found = true
				prevHead := member.IsHead
				prevManager := member.IsManager
				member.IsHead = isHead
				member.IsManager = isHead || (prevManager && !prevHead)
				next = append(next, member)
				continue
			}
			if isHead && member.IsHead {
				wasHead := true
				wasManager := member.IsManager
				member.IsHead = false
				member.IsManager = wasManager && !wasHead
			}
			next = append(next, member)
		}
		if keep && !found {
			next = append(next, entity.DepartmentMember{
				DepartmentID: deptID,
				UserID:       userID,
				IsHead:       isHead,
				IsManager:    isHead,
			})
		}
		if err := uc.deptMemberRepo.ReplaceByDepartment(deptID, next); err != nil {
			return nil, err
		}
		if isHead {
			if err := uc.deptRepo.UpdateHeadUserID(deptID, userID); err != nil {
				return nil, err
			}
		} else if dept.HeadUserID == userID {
			if err := uc.deptRepo.UpdateHeadUserID(deptID, 0); err != nil {
				return nil, err
			}
		}
	}

	existingRoles, err := uc.deptRoleRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	for _, role := range existingRoles {
		if err := uc.deptRoleRepo.Delete(role.ID); err != nil {
			return nil, err
		}
	}
	for _, row := range want {
		if err := uc.deptRoleRepo.Create(userID, row.departmentID, row.roleName); err != nil {
			return nil, err
		}
	}
	if err := uc.userDeptRepo.DeleteByUserID(userID); err != nil {
		return nil, err
	}
	for _, row := range want {
		if err := uc.userDeptRepo.Create(userID, row.departmentID); err != nil {
			return nil, err
		}
	}
	return allowed, nil
}

func (uc *UserAdmin) syncTeamLeaders(userID uint, teams []TeamAssignmentInput, allowedDepartments map[uint]struct{}) error {
	want := map[uint]struct{}{}
	for _, row := range teams {
		if row.TeamID == 0 {
			return fmt.Errorf("チームを選択してください")
		}
		if _, dup := want[row.TeamID]; dup {
			return fmt.Errorf("同じチームを重複して登録できません")
		}
		team, err := uc.teamRepo.GetByID(row.TeamID)
		if err != nil {
			return err
		}
		if _, ok := allowedDepartments[team.DepartmentID]; !ok {
			return fmt.Errorf("選択した部署に属さないチームは登録できません")
		}
		want[row.TeamID] = struct{}{}
	}

	existing, err := uc.teamMemberRepo.GetByMemberID(userID)
	if err != nil {
		return err
	}
	for _, member := range existing {
		if !member.IsReader {
			continue
		}
		if _, ok := want[member.TeamID]; ok {
			continue
		}
		if err := uc.teamMemberRepo.Delete(member.ID); err != nil {
			return err
		}
	}
	for teamID := range want {
		members, err := uc.teamMemberRepo.GetByTeamID(teamID)
		if err != nil {
			return err
		}
		found := false
		for i := range members {
			member := &members[i]
			isReader := member.MemberID == userID
			if isReader {
				found = true
			}
			if member.IsReader == isReader {
				continue
			}
			member.IsReader = isReader
			if err := uc.teamMemberRepo.Update(member.ID, member); err != nil {
				return err
			}
		}
		if !found {
			if err := uc.teamMemberRepo.Create(teamID, userID, true); err != nil {
				return err
			}
		}
	}
	return nil
}

func (uc *UserAdmin) syncOperationTargets(userID uint, target valueobject.OperationTarget, targetIDs []uint) error {
	if err := uc.deptRefRepo.DeleteByUserID(userID); err != nil {
		return err
	}
	if err := uc.teamRefRepo.DeleteByUserID(userID); err != nil {
		return err
	}
	if err := uc.relationshipRepo.DeleteByBossUserID(userID); err != nil {
		return err
	}
	switch target {
	case valueobject.OperationTargetDepartment:
		for _, id := range targetIDs {
			if err := uc.deptRefRepo.Create(userID, id); err != nil {
				return err
			}
		}
	case valueobject.OperationTargetTeam:
		for _, id := range targetIDs {
			if err := uc.teamRefRepo.Create(userID, id); err != nil {
				return err
			}
		}
	case valueobject.OperationTargetSubordinate:
		for _, subID := range targetIDs {
			if subID == userID {
				continue
			}
			if err := uc.relationshipRepo.Create(subID, userID, valueobject.RelationshipTypeBoss); err != nil {
				return err
			}
		}
	}
	return nil
}

func (uc *UserAdmin) SetUsersSuspended(ids []uint, suspended bool) error {
	if len(ids) == 0 {
		return fmt.Errorf("対象のユーザが選択されていません")
	}
	return uc.userRepo.SetSuspendedByIDs(ids, suspended)
}

func ValidateLineWorksLocalID(id string) error {
	if strings.Contains(id, "@") {
		return fmt.Errorf("LINE WORKS ID に @ を含めることはできません")
	}
	return nil
}
