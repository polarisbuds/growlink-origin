package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"
	"log/slog"

	"gorm.io/gorm"
)

type DepartmentRoleRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.DepartmentRole]
	unScopedBaseRepository BaseUnScopedRepository[entity.DepartmentRole]
}

func NewDepartmentRoleRepository(ctx context.Context, dp *gorm.DB) *DepartmentRoleRepository {
	return &DepartmentRoleRepository{
		baseRepository: BaseRepository[entity.DepartmentRole]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.DepartmentRole]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *DepartmentRoleRepository) Create(userID, departmentID uint, name string) error {
	if err := repo.baseRepository.Insert(&entity.DepartmentRole{
		Name:         name,
		UserID:       userID,
		DepartmentID: departmentID,
	}); err != nil {
		return fmt.Errorf("部署役職の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentRoleRepository) Update(id uint, role *entity.DepartmentRole) error {
	if role == nil {
		slog.Warn("部署役職の更新で空のDepartmentRoleが送られました")
		return nil
	}
	if _, err := repo.baseRepository.Updates(role, "id = ?", id); err != nil {
		return fmt.Errorf("部署役職の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentRoleRepository) Delete(id uint) error {
	if _, err := repo.baseRepository.Delete(1, "id = ?", id); err != nil {
		return fmt.Errorf("部署役職の削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentRoleRepository) GetAll() ([]entity.DepartmentRole, error) {
	roles, err := gorm.G[entity.DepartmentRole](repo.baseRepository.dp).Order("id ASC").
		Preload("User", nil).
		Preload("Department", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署役職の全取得に失敗しました: %w", err)
	}
	return roles, nil
}

func (repo *DepartmentRoleRepository) preloadWhere(query string, args ...any) ([]entity.DepartmentRole, error) {
	whereChain, err := repo.baseRepository.WhereChain("id ASC", query, args...)
	if err != nil {
		return nil, fmt.Errorf("部署役職情報取得用のチェインの作成に失敗しました: %w", err)
	}
	if whereChain == nil {
		return []entity.DepartmentRole{}, nil
	}
	roles, err := whereChain.Preload("User", nil).
		Preload("Department", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署役職情報の取得に失敗しました: %w", err)
	}
	return roles, nil
}

func (repo *DepartmentRoleRepository) GetByUserID(userID uint) ([]entity.DepartmentRole, error) {
	roles, err := repo.preloadWhere("user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザID指定での部署役職取得に失敗しました: %w", err)
	}
	return roles, nil
}

func (repo *DepartmentRoleRepository) GetByDepartmentID(departmentID uint) ([]entity.DepartmentRole, error) {
	roles, err := repo.preloadWhere("department_id = ?", departmentID)
	if err != nil {
		return nil, fmt.Errorf("部署ID指定での部署役職取得に失敗しました: %w", err)
	}
	return roles, nil
}
