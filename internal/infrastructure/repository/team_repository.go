package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type TeamRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Team]
	unScopedBaseRepository BaseUnScopedRepository[entity.Team]
}

func NewTeamsRepository(ctx context.Context, dp *gorm.DB) *TeamRepository {
	return &TeamRepository{
		baseRepository: BaseRepository[entity.Team]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Team]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *TeamRepository) GetByID(id uint) (*entity.Team, error) {
	team, err := repo.baseRepository.First("id ASC", "id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("チームの取得に失敗しました: %w", err)
	}
	if team == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return team, nil
}

func (repo *TeamRepository) GetAll() ([]entity.Team, error) {
	teams, err := gorm.G[entity.Team](repo.baseRepository.dp).Order("name ASC").
		Preload("Department", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("チームの全取得に失敗しました: %w", err)
	}
	return teams, nil
}

func (repo *TeamRepository) GetByDepartmentID(departmentID uint) ([]entity.Team, error) {
	teams, err := repo.baseRepository.Where("name ASC", "department_id = ?", departmentID)
	if err != nil {
		return nil, fmt.Errorf("部署ID指定でのチーム取得に失敗しました: %w", err)
	}
	return teams, nil
}

func (repo *TeamRepository) Create(team *entity.Team) error {
	if team == nil {
		return fmt.Errorf("作成対象のチームが空です")
	}
	if err := repo.baseRepository.Insert(team); err != nil {
		return fmt.Errorf("チームの作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamRepository) UpdateByID(team *entity.Team) error {
	if team == nil {
		return nil
	}
	if _, err := repo.baseRepository.Updates(team, "id = ?", team.ID); err != nil {
		return fmt.Errorf("チームの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamRepository) Delete(id uint) error {
	if _, err := repo.baseRepository.Delete(1, "id = ?", id); err != nil {
		return fmt.Errorf("チームの削除に失敗しました: %w", err)
	}
	return nil
}
