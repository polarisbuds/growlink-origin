package repository

import (
	"context"
	"errors"
	"fmt"
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DepartmentRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewDepartementRepository(ctx context.Context, dp *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{ctx: ctx, db: dp}
}

func (repo *DepartmentRepository) GetByID(id uint) (*entity.Department, error) {
	dept, err := gorm.G[entity.Department](repo.db).Where("id = ?", id).Order("id ASC").First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署の取得に失敗しました: %w", err)
	}
	return &dept, nil
}

func (repo *DepartmentRepository) GetByUUID(deptUUID uuid.UUID) (*entity.Department, error) {
	dept, err := gorm.G[entity.Department](repo.db).Where("uuid = ?", deptUUID.String()).Order("id ASC").First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("指定された部署UUIDは存在しませんでした: %w", err)
	}
	return &dept, nil
}

func (repo *DepartmentRepository) GetAll() ([]entity.Department, error) {
	depts, err := gorm.G[entity.Department](repo.db).Order("name ASC").Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署の全取得に失敗しました: %w", err)
	}
	return depts, nil
}

func (repo *DepartmentRepository) Create(department *entity.Department) error {
	if department == nil {
		return fmt.Errorf("作成対象の部署が空です")
	}
	if err := gorm.G[entity.Department](repo.db).Create(repo.ctx, department); err != nil {
		return fmt.Errorf("部署の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentRepository) UpdateByID(department *entity.Department) error {
	if department == nil {
		return nil
	}
	if _, err := gorm.G[entity.Department](repo.db).Where("id = ?", department.ID).Omit(clause.Associations).Updates(repo.ctx, *department); err != nil {
		return fmt.Errorf("部署の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentRepository) UpdateHeadUserID(id, headUserID uint) error {
	if err := repo.db.WithContext(repo.ctx).Model(&entity.Department{}).Where("id = ?", id).Update("head_user_id", headUserID).Error; err != nil {
		return fmt.Errorf("部署長の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentRepository) Delete(id uint) error {
	if _, err := gorm.G[entity.Department](repo.db).Where("id = ?", id).Delete(repo.ctx); err != nil {
		return fmt.Errorf("部署の削除に失敗しました: %w", err)
	}
	return nil
}
