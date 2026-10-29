package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type RoleRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewRoleRepository(ctx context.Context, dp *gorm.DB) *RoleRepository {
	return &RoleRepository{ctx: ctx, db: dp}
}

func (repo *RoleRepository) GetAll() ([]entity.Role, error) {
	roles, err := gorm.G[entity.Role](repo.db).Order("id ASC").Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ロールの全取得に失敗しました: %w", err)
	}
	return roles, nil
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
