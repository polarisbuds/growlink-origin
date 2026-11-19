package role_admin_usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/entity"
	logutils "growlink/internal/shared/logUtils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRoleNameEmpty   = errors.New("役職名を入力してください")
	ErrRoleNameInUse   = errors.New("同じ名前の役職が既に存在します")
	ErrRoleHasMembers  = errors.New("この役職はユーザに割り当てられているため削除できません")
)

type RoleOptionOutput struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type RoleAdmin struct {
	roleRepo repository_interface.RoleRepository
}

func NewRoleAdmin(roleRepo repository_interface.RoleRepository) *RoleAdmin {
	return &RoleAdmin{roleRepo: roleRepo}
}

func (uc *RoleAdmin) List() (ret0 []RoleOptionOutput, err error) {
	const op = "role_admin_usecase.List"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	roles, err := uc.roleRepo.GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]RoleOptionOutput, len(roles))
	for idx, r := range roles {
		out[idx] = RoleOptionOutput{UUID: r.UUID.String(), Name: r.Name}
	}
	return out, nil
}

func (uc *RoleAdmin) Create(name string) (ret0 *RoleOptionOutput, err error) {
	const op = "role_admin_usecase.Create"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrRoleNameEmpty
	}
	inUse, err := uc.nameInUse(name, 0)
	if err != nil {
		return nil, err
	}
	if inUse {
		return nil, ErrRoleNameInUse
	}
	role := &entity.Role{Name: name}
	if err := uc.roleRepo.Create(role); err != nil {
		return nil, err
	}
	return &RoleOptionOutput{UUID: role.UUID.String(), Name: role.Name}, nil
}

func (uc *RoleAdmin) Update(roleUUID uuid.UUID, name string) (ret0 *RoleOptionOutput, err error) {
	const op = "role_admin_usecase.Update"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrRoleNameEmpty
	}
	role, err := uc.roleRepo.GetByUUID(roleUUID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, gorm.ErrRecordNotFound
	}
	inUse, err := uc.nameInUse(name, role.ID)
	if err != nil {
		return nil, err
	}
	if inUse {
		return nil, ErrRoleNameInUse
	}
	if err := uc.roleRepo.Update(&entity.Role{Model: gorm.Model{ID: role.ID}, Name: name}); err != nil {
		return nil, err
	}
	updated, err := uc.roleRepo.GetByUUID(roleUUID)
	if err != nil || updated == nil {
		return &RoleOptionOutput{UUID: roleUUID.String(), Name: name}, nil
	}
	return &RoleOptionOutput{UUID: updated.UUID.String(), Name: updated.Name}, nil
}

func (uc *RoleAdmin) DeleteByUUID(roleUUID uuid.UUID) (err error) {
	const op = "role_admin_usecase.DeleteByUUID"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	role, err := uc.roleRepo.GetByUUID(roleUUID)
	if err != nil {
		return err
	}
	if role == nil {
		return gorm.ErrRecordNotFound
	}
	count, err := uc.roleRepo.CountMembersByRoleID(role.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrRoleHasMembers
	}
	return uc.roleRepo.Delete(role.ID)
}

func (uc *RoleAdmin) ResolveID(roleUUID uuid.UUID) (uint, error) {
	role, err := uc.roleRepo.GetByUUID(roleUUID)
	if err != nil {
		return 0, err
	}
	if role == nil {
		return 0, gorm.ErrRecordNotFound
	}
	return role.ID, nil
}

func (uc *RoleAdmin) nameInUse(name string, excludeID uint) (bool, error) {
	existing, err := uc.roleRepo.GetByName(name)
	if err != nil {
		return false, fmt.Errorf("役職名の重複確認に失敗しました: %w", err)
	}
	if existing == nil {
		return false, nil
	}
	return existing.ID != excludeID, nil
}