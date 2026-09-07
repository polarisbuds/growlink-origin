package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type UserReferenceUserRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.UserReferenceUser]
	unScopedBaseRepository BaseUnScopedRepository[entity.UserReferenceUser]
}

func NewUserReferenceUserRepository(ctx context.Context, dp *gorm.DB) *UserReferenceUserRepository {
	return &UserReferenceUserRepository{
		baseRepository: BaseRepository[entity.UserReferenceUser]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.UserReferenceUser]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *UserReferenceUserRepository) Create(userID, referableUserID uint) error {
	if err := repo.baseRepository.Insert(&entity.UserReferenceUser{
		UserID:          userID,
		ReferableUserID: referableUserID,
	}); err != nil {
		return fmt.Errorf("ユーザ参照（ユーザ）の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserReferenceUserRepository) Delete(id uint) error {
	if _, err := repo.baseRepository.Delete(1, "id = ?", id); err != nil {
		return fmt.Errorf("ユーザ参照（ユーザ）の削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserReferenceUserRepository) GetAll() ([]entity.UserReferenceUser, error) {
	references, err := gorm.G[entity.UserReferenceUser](repo.baseRepository.dp).Order("id ASC").
		Preload("User", nil).
		Preload("ReferableUser", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ユーザ参照（ユーザ）の全取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceUserRepository) preloadWhere(query string, args ...any) ([]entity.UserReferenceUser, error) {
	whereChain, err := repo.baseRepository.WhereChain("id ASC", query, args...)
	if err != nil {
		return nil, fmt.Errorf("ユーザ参照（ユーザ）情報取得用のチェインの作成に失敗しました: %w", err)
	}
	if whereChain == nil {
		return []entity.UserReferenceUser{}, nil
	}
	references, err := whereChain.Preload("User", nil).
		Preload("ReferableUser", nil).Find(repo.ctx)
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
