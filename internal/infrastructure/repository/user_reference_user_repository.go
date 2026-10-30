package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type UserReferenceUserRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewUserReferenceUserRepository(ctx context.Context, dp *gorm.DB) *UserReferenceUserRepository {
	return &UserReferenceUserRepository{ctx: ctx, db: dp}
}

func (repo *UserReferenceUserRepository) Create(userID, referableUserID uint) error {
	if err := gorm.G[entity.UserReferenceUser](repo.db).Create(repo.ctx, &entity.UserReferenceUser{
		UserID:          userID,
		ReferableUserID: referableUserID,
	}); err != nil {
		return fmt.Errorf("ユーザ参照（ユーザ）の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserReferenceUserRepository) Delete(id uint) error {
	if _, err := gorm.G[entity.UserReferenceUser](repo.db).Where("id = ?", id).Delete(repo.ctx); err != nil {
		return fmt.Errorf("ユーザ参照（ユーザ）の削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserReferenceUserRepository) GetAll() ([]entity.UserReferenceUser, error) {
	references, err := gorm.G[entity.UserReferenceUser](repo.db).Order("id ASC").
		Preload("User", nil).
		Preload("ReferableUser", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ユーザ参照（ユーザ）の全取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceUserRepository) preloadWhere(query string, args ...any) ([]entity.UserReferenceUser, error) {
	references, err := gorm.G[entity.UserReferenceUser](repo.db).
		Where(query, args...).
		Order("id ASC").
		Preload("User", nil).
		Preload("ReferableUser", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ユーザ参照（ユーザ）情報の取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceUserRepository) GetByUserID(userID uint) ([]entity.UserReferenceUser, error) {
	references, err := repo.preloadWhere("user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザID指定でのユーザ参照（ユーザ）取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceUserRepository) GetByReferableUserID(referableUserID uint) ([]entity.UserReferenceUser, error) {
	references, err := repo.preloadWhere("referable_user_id = ?", referableUserID)
	if err != nil {
		return nil, fmt.Errorf("参照可能ユーザID指定でのユーザ参照（ユーザ）取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceUserRepository) IsReferable(userID, referableUserID uint) (isReferable bool, err error) {
	references, err := repo.preloadWhere("user_id = ? AND referable_user_id = ?", userID, referableUserID)
	if err != nil {
		return false, fmt.Errorf("ユーザーIDと参照可能なユーザID指定での参照権限有無のチェックに失敗しました: %w", err)
	}
	if len(references) > 0 {
		return true, nil
	}
	return false, nil
}
