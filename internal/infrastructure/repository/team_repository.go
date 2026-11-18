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

type TeamRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewTeamsRepository(ctx context.Context, dp *gorm.DB) *TeamRepository {
	return &TeamRepository{ctx: ctx, db: dp}
}

func (repo *TeamRepository) GetByID(id uint) (ret0 *entity.Team, err error) {
	team, err := gorm.G[entity.Team](repo.db).Where("id = ?", id).Order("id ASC").First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("チームの取得に失敗しました: %w", err)
	}
	return &team, nil
}

func (repo *TeamRepository) GetByIDs(ids []uint) (ret0 []entity.Team, err error) {
	team, err := gorm.G[entity.Team](repo.db).Where("id IN ?", ids).Order("id ASC").Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("チームの取得に失敗しました: %w", err)
	}
	return team, nil
}

func (repo *TeamRepository) GetByUUID(teamUUID uuid.UUID) (ret0 *entity.Team, err error) {
	team, err := gorm.G[entity.Team](repo.db).Where("uuid = ?", teamUUID.String()).Order("id ASC").First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("指定されたチームUUIDは存在しませんでした: %w", err)
	}
	return &team, nil
}

func (repo *TeamRepository) GetByUUIDs(teamUUIDs []string) (ret0 []entity.Team, err error) {
	teams, err := gorm.G[entity.Team](repo.db).Where("uuid IN ?", teamUUIDs).
		Preload("TeamMembers", nil).
		Preload("TeamMembers.Member", nil).
		Order("id ASC").Find(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("指定されたチームUUIDは存在しませんでした: %w", err)
	}
	return teams, nil
}

func (repo *TeamRepository) GetAll() (ret0 []entity.Team, err error) {
	teams, err := gorm.G[entity.Team](repo.db).Order("name ASC").
		Preload("Department", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("チームの全取得に失敗しました: %w", err)
	}
	return teams, nil
}

func (repo *TeamRepository) GetByDepartmentID(departmentID uint) (ret0 []entity.Team, err error) {
	teams, err := gorm.G[entity.Team](repo.db).
		Where("department_id = ?", departmentID).
		Order("name ASC").
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署ID指定でのチーム取得に失敗しました: %w", err)
	}
	return teams, nil
}

func (repo *TeamRepository) Create(team *entity.Team) (err error) {
	if team == nil {
		return fmt.Errorf("作成対象のチームが空です")
	}
	if err = gorm.G[entity.Team](repo.db).Create(repo.ctx, team); err != nil {
		return fmt.Errorf("チームの作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamRepository) UpdateByID(team *entity.Team) (err error) {
	if team == nil {
		return nil
	}
	if _, err = gorm.G[entity.Team](repo.db).Where("id = ?", team.ID).Omit(clause.Associations).Updates(repo.ctx, *team); err != nil {
		return fmt.Errorf("チームの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamRepository) Delete(id uint) (err error) {
	if _, err = gorm.G[entity.Team](repo.db).Where("id = ?", id).Delete(repo.ctx); err != nil {
		return fmt.Errorf("チームの削除に失敗しました: %w", err)
	}
	return nil
}
