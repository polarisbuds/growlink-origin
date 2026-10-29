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
	userRepo       repository_interface.UserRepository
	deptRefRepo    repository_interface.UserReferenceDepartmentRepository
	teamRefRepo    repository_interface.UserReferenceTeamRepository
	relationshipRepo repository_interface.UserRelationshipRepository
}

func NewUserAdmin(
	userRepo repository_interface.UserRepository,
	deptRefRepo repository_interface.UserReferenceDepartmentRepository,
	teamRefRepo repository_interface.UserReferenceTeamRepository,
	relationshipRepo repository_interface.UserRelationshipRepository,
) *UserAdmin {
	return &UserAdmin{
		userRepo:         userRepo,
		deptRefRepo:      deptRefRepo,
		teamRefRepo:      teamRefRepo,
		relationshipRepo: relationshipRepo,
	}
}

func (uc *UserAdmin) ListUsers() ([]outputmodel.ManagementUserOutput, error) {
	users, err := uc.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]outputmodel.ManagementUserOutput, 0, len(users))
	for _, u := range users {
		targetIDs, err := uc.loadTargetIDs(u.ID, u.OperationTarget)
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
			Suspended:       u.Suspended,
		})
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

type CreateUserInput struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	DepartmentID uint   `json:"departmentId"`
	RoleID       uint   `json:"roleId"`
	Password     string `json:"password"`
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
	if input.DepartmentID == 0 {
		return nil, fmt.Errorf("所属部署を選択してください")
	}
	if input.RoleID == 0 {
		return nil, fmt.Errorf("役職を選択してください")
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
		DepartmentID:    input.DepartmentID,
		RoleID:          input.RoleID,
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
		Suspended:       user.Suspended,
	}, nil
}

type UpdateUserInput struct {
	Authority       string `json:"authority"`
	SystemRole      string `json:"systemRole"`
	OperationTarget string `json:"operationTarget"`
	TargetIDs       []uint `json:"targetIds"`
	DepartmentID    uint   `json:"departmentId"`
	RoleID          uint   `json:"roleId"`
	Suspended       *bool  `json:"suspended"`
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
	if input.DepartmentID > 0 {
		user.DepartmentID = input.DepartmentID
	}
	if input.RoleID > 0 {
		user.RoleID = input.RoleID
	}
	if input.Suspended != nil {
		user.Suspended = *input.Suspended
	}
	if err := uc.userRepo.UpdateByID(user); err != nil {
		return err
	}
	return uc.syncOperationTargets(userID, user.OperationTarget, input.TargetIDs)
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
