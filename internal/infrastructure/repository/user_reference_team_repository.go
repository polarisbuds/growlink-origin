package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type UserReferenceTeamRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewUserReferenceTeamRepository(ctx context.Context, dp *gorm.DB) *UserReferenceTeamRepository {
	return &UserReferenceTeamRepository{ctx: ctx, db: dp}
}

func (repo *UserReferenceTeamRepository) Create(userID, referableTeamID uint) error {
	if err := gorm.G[entity.UserReferenceTeam](repo.db).Create(repo.ctx, &entity.UserReferenceTeam{
		UserID:          userID,
		ReferableTeamID: referableTeamID,
	}); err != nil {
		return fmt.Errorf("ユーザ参照（チーム）の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserReferenceTeamRepository) Delete(id uint) error {
	if _, err := gorm.G[entity.UserReferenceTeam](repo.db).Where("id = ?", id).Delete(repo.ctx); err != nil {
		return fmt.Errorf("ユーザ参照（チーム）の削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserReferenceTeamRepository) GetAll() ([]entity.UserReferenceTeam, error) {
	references, err := gorm.G[entity.UserReferenceTeam](repo.db).Order("id ASC").
		Preload("User", nil).
		Preload("ReferableTeam", nil).
		Preload("ReferableTeam.Department", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ユーザ参照（チーム）の全取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceTeamRepository) preloadWhere(query string, args ...any) ([]entity.UserReferenceTeam, error) {
	references, err := gorm.G[entity.UserReferenceTeam](repo.db).
		Where(query, args...).
		Order("id ASC").
		Preload("User", nil).
		Preload("ReferableTeam", nil).
		Preload("ReferableTeam.Department", nil).
		Find(repo.ctx)
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
	if _, err := gorm.G[entity.UserReferenceTeam](repo.db).Where("user_id = ?", userID).Delete(repo.ctx); err != nil {
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
