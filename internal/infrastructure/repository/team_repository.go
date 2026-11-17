package repository

import (
	"context"
	"errors"
	"fmt"
	"growlink/internal/domain/entity"
	logutils "growlink/internal/shared/logUtils"

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
	const op = "team_repository.GetByID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	team, err := gorm.G[entity.Team](repo.db).Where("id = ?", id).Order("id ASC").First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("チームの取得に失敗しました: %w", err)
	}
	return &team, nil
}

func (repo *TeamRepository) GetByUUID(teamUUID uuid.UUID) (ret0 *entity.Team, err error) {
	const op = "team_repository.GetByUUID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	team, err := gorm.G[entity.Team](repo.db).Where("uuid = ?", teamUUID.String()).Order("id ASC").First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("指定されたチームUUIDは存在しませんでした: %w", err)
	}
	return &team, nil
}

func (repo *TeamRepository) GetAll() (ret0 []entity.Team, err error) {
	const op = "team_repository.GetAll"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	teams, err := gorm.G[entity.Team](repo.db).Order("name ASC").
		Preload("Department", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("チームの全取得に失敗しました: %w", err)
	}
	return teams, nil
}

func (repo *TeamRepository) GetByDepartmentID(departmentID uint) (ret0 []entity.Team, err error) {
	const op = "team_repository.GetByDepartmentID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
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
	const op = "team_repository.Create"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if team == nil {
		return fmt.Errorf("作成対象のチームが空です")
	}
	if err = gorm.G[entity.Team](repo.db).Create(repo.ctx, team); err != nil {
		return fmt.Errorf("チームの作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamRepository) UpdateByID(team *entity.Team) (err error) {
	const op = "team_repository.UpdateByID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if team == nil {
		return nil
	}
	if _, err = gorm.G[entity.Team](repo.db).Where("id = ?", team.ID).Omit(clause.Associations).Updates(repo.ctx, *team); err != nil {
		return fmt.Errorf("チームの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamRepository) Delete(id uint) (err error) {
	const op = "team_repository.Delete"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if _, err = gorm.G[entity.Team](repo.db).Where("id = ?", id).Delete(repo.ctx); err != nil {
		return fmt.Errorf("チームの削除に失敗しました: %w", err)
	}
	return nil
}
