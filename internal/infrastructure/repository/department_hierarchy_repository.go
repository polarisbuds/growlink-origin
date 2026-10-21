package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type DepartmentHierarchyRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.DepartmentHierarchy]
	unScopedBaseRepository BaseUnScopedRepository[entity.DepartmentHierarchy]
}

func NewDepartmentHierarchyRepository(ctx context.Context, dp *gorm.DB) *DepartmentHierarchyRepository {
	return &DepartmentHierarchyRepository{
		baseRepository: BaseRepository[entity.DepartmentHierarchy]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.DepartmentHierarchy]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *DepartmentHierarchyRepository) Create(parentDepartmentID, childDepartmentID uint) error {
	if err := repo.baseRepository.Insert(&entity.DepartmentHierarchy{
		ParentDepartmentID: parentDepartmentID,
		ChildDepartmentID:  childDepartmentID,
	}); err != nil {
		return fmt.Errorf("部署階層の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentHierarchyRepository) Delete(id uint) error {
	if _, err := repo.baseRepository.Delete(1, "id = ?", id); err != nil {
		return fmt.Errorf("部署階層の削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentHierarchyRepository) GetAll() ([]entity.DepartmentHierarchy, error) {
	hierarchies, err := gorm.G[entity.DepartmentHierarchy](repo.baseRepository.dp).Order("id ASC").
		Preload("ParentDepartment", nil).
		Preload("ChildDepartment", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署階層の全取得に失敗しました: %w", err)
	}
	return hierarchies, nil
}

func (repo *DepartmentHierarchyRepository) preloadWhere(query string, args ...any) ([]entity.DepartmentHierarchy, error) {
	whereChain, err := repo.baseRepository.WhereChain("id ASC", query, args...)
	if err != nil {
		return nil, fmt.Errorf("部署階層情報取得用のチェインの作成に失敗しました: %w", err)
	}
	if whereChain == nil {
		return []entity.DepartmentHierarchy{}, nil
	}
	hierarchies, err := whereChain.Preload("ParentDepartment", nil).
		Preload("ChildDepartment", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署階層情報の取得に失敗しました: %w", err)
	}
	return hierarchies, nil
}

func (repo *DepartmentHierarchyRepository) GetByParentDepartmentID(parentDepartmentID uint) ([]entity.DepartmentHierarchy, error) {
	hierarchies, err := repo.preloadWhere("parent_department_id = ?", parentDepartmentID)
	if err != nil {
		return nil, fmt.Errorf("親部署ID指定での部署階層取得に失敗しました: %w", err)
	}
	return hierarchies, nil
}

func (repo *DepartmentHierarchyRepository) GetByChildDepartmentID(childDepartmentID uint) ([]entity.DepartmentHierarchy, error) {
	hierarchies, err := repo.preloadWhere("child_department_id = ?", childDepartmentID)
	if err != nil {
		return nil, fmt.Errorf("子部署ID指定での部署階層取得に失敗しました: %w", err)
	}
	return hierarchies, nil
}

func (repo *DepartmentHierarchyRepository) DeleteByParentDepartmentID(parentDepartmentID uint) error {
	if _, err := repo.baseRepository.Delete(0, "parent_department_id = ?", parentDepartmentID); err != nil {
		return fmt.Errorf("親部署ID指定での部署階層削除に失敗しました: %w", err)
	}
	return nil
}
