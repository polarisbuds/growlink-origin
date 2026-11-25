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
	ctx          context.Context
	db, unscoped *gorm.DB
}

func NewDepartementRepository(ctx context.Context, dp *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{ctx: ctx, db: dp, unscoped: dp.Session(&gorm.Session{PropagateUnscoped: true})}
}

func (repo *DepartmentRepository) GetByID(id uint) (ret0 *entity.Department, err error) {
	dept, err := gorm.G[entity.Department](repo.db).Where("id = ?", id).Order("id ASC").First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署の取得に失敗しました: %w", err)
	}
	return &dept, nil
}

func (repo *DepartmentRepository) GetByIDs(ids []uint) (ret0 []entity.Department, err error) {
	depts, err := gorm.G[entity.Department](repo.db).Where("id IN ?", ids).
		Preload("Teams", nil).
		Preload("Teams.TeamMembers", nil).
		Preload("Teams.TeamMembers.Member", nil).
		Preload("DepartmentMembers", nil).
		Preload("DepartmentMembers.User", nil).
		Order("id ASC").Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署の取得に失敗しました: %w", err)
	}
	return depts, nil
}

func (repo *DepartmentRepository) GetByIDsByUnscoped(ids []uint) (ret0 []entity.Department, err error) {
	depts, err := gorm.G[entity.Department](repo.unscoped).Where("id IN ?", ids).
		Scopes(func(stmt *gorm.Statement) { stmt.Unscoped = true }).
		Preload("Teams", nil).
		Preload("Teams.TeamMembers", nil).
		Preload("Teams.TeamMembers.Member", nil).
		Preload("DepartmentMembers", nil).
		Preload("DepartmentMembers.User", nil).
		Order("id ASC").Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署の取得に失敗しました: %w", err)
	}
	return depts, nil
}

func (repo *DepartmentRepository) GetByName(name string) (ret0 *entity.Department, err error) {
	if cnt, err := gorm.G[entity.Department](repo.db).Where("name = ?", name).Count(repo.ctx, "id"); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("指定された部署名は存在しませんでした: %w", err)
	} else if cnt == 0 {
		return nil, nil
	}

	dept, err := gorm.G[entity.Department](repo.db).Where("name = ?", name).Order("id ASC").First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("指定された部署名は存在しませんでした: %w", err)
	}
	return &dept, nil
}

func (repo *DepartmentRepository) GetByUUID(deptUUID uuid.UUID) (ret0 *entity.Department, err error) {
	dept, err := gorm.G[entity.Department](repo.db).Where("uuid = ?", deptUUID.String()).Order("id ASC").First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("指定された部署UUIDは存在しませんでした: %w", err)
	}
	return &dept, nil
}

func (repo *DepartmentRepository) GetByUUIDs(deptUUIDs []string) (ret0 []entity.Department, err error) {
	depts, err := gorm.G[entity.Department](repo.db).Where("uuid IN ?", deptUUIDs).
		Preload("Teams", nil).
		Preload("Teams.TeamMembers", nil).
		Preload("Teams.TeamMembers.Member", nil).
		Preload("DepartmentMembers", nil).
		Preload("DepartmentMembers.User", nil).
		Order("id ASC").Find(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("指定された部署UUIDは存在しませんでした: %w", err)
	}
	return depts, nil
}

func (repo *DepartmentRepository) GetAll() (ret0 []entity.Department, err error) {
	depts, err := gorm.G[entity.Department](repo.db).Order("name ASC").Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署の全取得に失敗しました: %w", err)
	}
	return depts, nil
}

func (repo *DepartmentRepository) Create(department *entity.Department) (err error) {
	if department == nil {
		return fmt.Errorf("作成対象の部署が空です")
	}
	if err = gorm.G[entity.Department](repo.db).Create(repo.ctx, department); err != nil {
		return fmt.Errorf("部署の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentRepository) UpdateByID(department *entity.Department) (err error) {
	if department == nil {
		return nil
	}
	if _, err = gorm.G[entity.Department](repo.db).Where("id = ?", department.ID).Omit(clause.Associations).Updates(repo.ctx, *department); err != nil {
		return fmt.Errorf("部署の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentRepository) UpdateHeadUserID(id, headUserID uint) (err error) {
	if err = repo.db.WithContext(repo.ctx).Model(&entity.Department{}).Where("id = ?", id).Update("head_user_id", headUserID).Error; err != nil {
		return fmt.Errorf("部署長の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentRepository) Delete(id uint) (err error) {
	if _, err = gorm.G[entity.Department](repo.db).Where("id = ?", id).Delete(repo.ctx); err != nil {
		return fmt.Errorf("部署の削除に失敗しました: %w", err)
	}
	return nil
}
