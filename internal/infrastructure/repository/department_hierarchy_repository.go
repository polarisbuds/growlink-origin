package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"
	logutils "growlink/internal/shared/logUtils"

	"gorm.io/gorm"
)

type DepartmentHierarchyRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewDepartmentHierarchyRepository(ctx context.Context, dp *gorm.DB) *DepartmentHierarchyRepository {
	return &DepartmentHierarchyRepository{ctx: ctx, db: dp}
}

func (repo *DepartmentHierarchyRepository) Create(parentDepartmentID, childDepartmentID uint) (err error) {
	const op = "department_hierarchy_repository.Create"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if err = gorm.G[entity.DepartmentHierarchy](repo.db).Create(repo.ctx, &entity.DepartmentHierarchy{
		ParentDepartmentID: parentDepartmentID,
		ChildDepartmentID:  childDepartmentID,
	}); err != nil {
		return fmt.Errorf("部署階層の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentHierarchyRepository) Update(id uint, hierarchy entity.DepartmentHierarchy) (err error) {
	const op = "department_hierarchy_repository.Update"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if _, err = gorm.G[entity.DepartmentHierarchy](repo.db).Updates(repo.ctx, hierarchy); err != nil {
		return fmt.Errorf("部署階層の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentHierarchyRepository) Delete(id uint) (err error) {
	const op = "department_hierarchy_repository.Delete"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if _, err = gorm.G[entity.DepartmentHierarchy](repo.db).Where("id = ?", id).Delete(repo.ctx); err != nil {
		return fmt.Errorf("部署階層の削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentHierarchyRepository) GetAll() (ret0 []entity.DepartmentHierarchy, err error) {
	const op = "department_hierarchy_repository.GetAll"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	hierarchies, err := gorm.G[entity.DepartmentHierarchy](repo.db).Order("id ASC").
		Preload("ParentDepartment", nil).
		Preload("ChildDepartment", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署階層の全取得に失敗しました: %w", err)
	}
	return hierarchies, nil
}

func (repo *DepartmentHierarchyRepository) preloadWhere(query string, args ...any) ([]entity.DepartmentHierarchy, error) {
	hierarchies, err := gorm.G[entity.DepartmentHierarchy](repo.db).
		Where(query, args...).
		Order("id ASC").
		Preload("ParentDepartment", nil).
		Preload("ChildDepartment", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署階層情報の取得に失敗しました: %w", err)
	}
	return hierarchies, nil
}

func (repo *DepartmentHierarchyRepository) GetByParentDepartmentID(parentDepartmentID uint) (ret0 []entity.DepartmentHierarchy, err error) {
	const op = "department_hierarchy_repository.GetByParentDepartmentID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	hierarchies, err := repo.preloadWhere("parent_department_id = ?", parentDepartmentID)
	if err != nil {
		return nil, fmt.Errorf("親部署ID指定での部署階層取得に失敗しました: %w", err)
	}
	return hierarchies, nil
}

func (repo *DepartmentHierarchyRepository) GetByChildDepartmentID(childDepartmentID uint) (ret0 []entity.DepartmentHierarchy, err error) {
	const op = "department_hierarchy_repository.GetByChildDepartmentID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	hierarchies, err := repo.preloadWhere("child_department_id = ?", childDepartmentID)
	if err != nil {
		return nil, fmt.Errorf("子部署ID指定での部署階層取得に失敗しました: %w", err)
	}
	return hierarchies, nil
}

func (repo *DepartmentHierarchyRepository) DeleteByParentDepartmentID(parentDepartmentID uint) (err error) {
	const op = "department_hierarchy_repository.DeleteByParentDepartmentID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if _, err = gorm.G[entity.DepartmentHierarchy](repo.db).
		Where("parent_department_id = ?", parentDepartmentID).
		Delete(repo.ctx); err != nil {
		return fmt.Errorf("親部署ID指定での部署階層削除に失敗しました: %w", err)
	}
	return nil
}
