package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Department]
	unScopedBaseRepository BaseUnScopedRepository[entity.Department]
}

func NewDepartementRepository(ctx context.Context, dp *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{
		baseRepository: BaseRepository[entity.Department]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Department]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *DepartmentRepository) GetByID(id uint) (*entity.Department, error) {
	dept, err := repo.baseRepository.First("id ASC", "id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("部署の取得に失敗しました: %w", err)
	}
	if dept == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return dept, nil
}

func (repo *DepartmentRepository) GetAll() ([]entity.Department, error) {
	depts, err := repo.baseRepository.All("name ASC")
	if err != nil {
		return nil, fmt.Errorf("部署の全取得に失敗しました: %w", err)
	}
	return depts, nil
}

func (repo *DepartmentRepository) Create(department *entity.Department) error {
	if department == nil {
		return fmt.Errorf("作成対象の部署が空です")
	}
	if err := repo.baseRepository.Insert(department); err != nil {
		return fmt.Errorf("部署の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentRepository) UpdateByID(department *entity.Department) error {
	if department == nil {
		return nil
	}
	if _, err := repo.baseRepository.Updates(department, "id = ?", department.ID); err != nil {
		return fmt.Errorf("部署の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentRepository) Delete(id uint) error {
	if _, err := repo.baseRepository.Delete(1, "id = ?", id); err != nil {
		return fmt.Errorf("部署の削除に失敗しました: %w", err)
	}
	return nil
}
