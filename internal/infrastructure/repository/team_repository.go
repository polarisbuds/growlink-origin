package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TeamRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewTeamsRepository(ctx context.Context, dp *gorm.DB) *TeamRepository {
	return &TeamRepository{ctx: ctx, db: dp}
}

func (repo *TeamRepository) GetByID(id uint) (*entity.Team, error) {
	team, err := gorm.G[entity.Team](repo.db).Where("id = ?", id).Order("id ASC").First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("チームの取得に失敗しました: %w", err)
	}
	return &team, nil
}

func (repo *TeamRepository) GetAll() ([]entity.Team, error) {
	teams, err := gorm.G[entity.Team](repo.db).Order("name ASC").
		Preload("Department", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("チームの全取得に失敗しました: %w", err)
	}
	return teams, nil
}

func (repo *TeamRepository) GetByDepartmentID(departmentID uint) ([]entity.Team, error) {
	teams, err := gorm.G[entity.Team](repo.db).
		Where("department_id = ?", departmentID).
		Order("name ASC").
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署ID指定でのチーム取得に失敗しました: %w", err)
	}
	return teams, nil
}

func (repo *TeamRepository) Create(team *entity.Team) error {
	if team == nil {
		return fmt.Errorf("作成対象のチームが空です")
	}
	if err := gorm.G[entity.Team](repo.db).Create(repo.ctx, team); err != nil {
		return fmt.Errorf("チームの作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamRepository) UpdateByID(team *entity.Team) error {
	if team == nil {
		return nil
	}
	if _, err := gorm.G[entity.Team](repo.db).Where("id = ?", team.ID).Omit(clause.Associations).Updates(repo.ctx, *team); err != nil {
		return fmt.Errorf("チームの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamRepository) Delete(id uint) error {
	if _, err := gorm.G[entity.Team](repo.db).Where("id = ?", id).Delete(repo.ctx); err != nil {
		return fmt.Errorf("チームの削除に失敗しました: %w", err)
	}
	return nil
}
