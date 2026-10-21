package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type UserReferenceTeamRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.UserReferenceTeam]
	unScopedBaseRepository BaseUnScopedRepository[entity.UserReferenceTeam]
}

func NewUserReferenceTeamRepository(ctx context.Context, dp *gorm.DB) *UserReferenceTeamRepository {
	return &UserReferenceTeamRepository{
		baseRepository: BaseRepository[entity.UserReferenceTeam]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.UserReferenceTeam]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *UserReferenceTeamRepository) Create(userID, referableTeamID uint) error {
	if err := repo.baseRepository.Insert(&entity.UserReferenceTeam{
		UserID:          userID,
		ReferableTeamID: referableTeamID,
	}); err != nil {
		return fmt.Errorf("ユーザ参照（チーム）の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserReferenceTeamRepository) Delete(id uint) error {
	if _, err := repo.baseRepository.Delete(1, "id = ?", id); err != nil {
		return fmt.Errorf("ユーザ参照（チーム）の削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserReferenceTeamRepository) GetAll() ([]entity.UserReferenceTeam, error) {
	references, err := gorm.G[entity.UserReferenceTeam](repo.baseRepository.dp).Order("id ASC").
		Preload("User", nil).
		Preload("ReferableTeam", nil).
		Preload("ReferableTeam.Department", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ユーザ参照（チーム）の全取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceTeamRepository) preloadWhere(query string, args ...any) ([]entity.UserReferenceTeam, error) {
	whereChain, err := repo.baseRepository.WhereChain("id ASC", query, args...)
	if err != nil {
		return nil, fmt.Errorf("ユーザ参照（チーム）情報取得用のチェインの作成に失敗しました: %w", err)
	}
	if whereChain == nil {
		return []entity.UserReferenceTeam{}, nil
	}
	references, err := whereChain.Preload("User", nil).
		Preload("ReferableTeam", nil).
		Preload("ReferableTeam.Department", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ユーザ参照（チーム）情報の取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceTeamRepository) GetByUserID(userID uint) ([]entity.UserReferenceTeam, error) {
	references, err := repo.preloadWhere("user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザID指定でのユーザ参照（チーム）取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceTeamRepository) GetByReferableTeamID(referableTeamID uint) ([]entity.UserReferenceTeam, error) {
	references, err := repo.preloadWhere("referable_team_id = ?", referableTeamID)
	if err != nil {
		return nil, fmt.Errorf("参照可能チームID指定でのユーザ参照（チーム）取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceTeamRepository) DeleteByUserID(userID uint) error {
	if _, err := repo.baseRepository.Delete(0, "user_id = ?", userID); err != nil {
		return fmt.Errorf("ユーザ参照（チーム）のユーザ単位削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserReferenceTeamRepository) IsReferable(userID, referableTeamID uint) (isReferable bool, err error) {
	references, err := repo.preloadWhere("user_id = ? AND referable_team_id = ?", userID, referableTeamID)
	if err != nil {
		return false, fmt.Errorf("ユーザーIDと参照可能なチームID指定での参照権限有無のチェックに失敗しました: %w", err)
	}
	if len(references) > 0 {
		return true, nil
	}
	return false, nil
}
