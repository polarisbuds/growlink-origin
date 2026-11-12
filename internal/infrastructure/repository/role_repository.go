package repository

import (
	"context"
	"errors"
	"fmt"
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewRoleRepository(ctx context.Context, dp *gorm.DB) *RoleRepository {
	return &RoleRepository{ctx: ctx, db: dp}
}

func (repo *RoleRepository) GetByUUID(roleUUID uuid.UUID) (*entity.Role, error) {
	role, err := gorm.G[entity.Role](repo.db).Where("uuid = ?", roleUUID).First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("UUID指定での役職取得に失敗しました: %w", err)
	}
	return &role, nil
}

func (repo *RoleRepository) GetAll() ([]entity.Role, error) {
	roles, err := gorm.G[entity.Role](repo.db).Order("id ASC").Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ロールの全取得に失敗しました: %w", err)
	}
	return roles, nil
}

func (repo *RoleRepository) GetByName(name string) (*entity.Role, error) {
	role, err := gorm.G[entity.Role](repo.db).Where("name = ?", name).First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("役職名での取得に失敗しました: %w", err)
	}
	return &role, nil
}

func (repo *RoleRepository) Create(role *entity.Role) error {
	if role == nil {
		return fmt.Errorf("役職が空です")
	}
	if err := gorm.G[entity.Role](repo.db).Create(repo.ctx, role); err != nil {
		return fmt.Errorf("役職の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *RoleRepository) Update(role *entity.Role) error {
	if role == nil || role.ID == 0 {
		return fmt.Errorf("役職が空です")
	}
	if _, err := gorm.G[entity.Role](repo.db).Where("id = ?", role.ID).Updates(repo.ctx, *role); err != nil {
		return fmt.Errorf("役職の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *RoleRepository) Delete(id uint) error {
	if _, err := gorm.G[entity.Role](repo.db).Where("id = ?", id).Delete(repo.ctx); err != nil {
		return fmt.Errorf("役職の削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *RoleRepository) CountMembersByRoleID(roleID uint) (int64, error) {
	var count int64
	if err := repo.db.WithContext(repo.ctx).Model(&entity.DepartmentMember{}).Where("role_id = ?", roleID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("役職の割当件数取得に失敗しました: %w", err)
	}
	return count, nil
}
