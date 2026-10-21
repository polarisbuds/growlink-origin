package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type UserReferenceDepartmentRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.UserReferenceDepartment]
	unScopedBaseRepository BaseUnScopedRepository[entity.UserReferenceDepartment]
}

func NewUserReferenceDepartmentRepository(ctx context.Context, dp *gorm.DB) *UserReferenceDepartmentRepository {
	return &UserReferenceDepartmentRepository{
		baseRepository: BaseRepository[entity.UserReferenceDepartment]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.UserReferenceDepartment]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *UserReferenceDepartmentRepository) Create(userID, referableDepartmentID uint) error {
	if err := repo.baseRepository.Insert(&entity.UserReferenceDepartment{
		UserID:                userID,
		ReferableDepartmentID: referableDepartmentID,
	}); err != nil {
		return fmt.Errorf("ユーザ参照（部署）の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserReferenceDepartmentRepository) Delete(id uint) error {
	if _, err := repo.baseRepository.Delete(1, "id = ?", id); err != nil {
		return fmt.Errorf("ユーザ参照（部署）の削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserReferenceDepartmentRepository) GetAll() ([]entity.UserReferenceDepartment, error) {
	references, err := gorm.G[entity.UserReferenceDepartment](repo.baseRepository.dp).Order("id ASC").
		Preload("User", nil).
		Preload("ReferableDepartment", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ユーザ参照（部署）の全取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceDepartmentRepository) preloadWhere(query string, args ...any) ([]entity.UserReferenceDepartment, error) {
	whereChain, err := repo.baseRepository.WhereChain("id ASC", query, args...)
	if err != nil {
		return nil, fmt.Errorf("ユーザ参照（部署）情報取得用のチェインの作成に失敗しました: %w", err)
	}
	if whereChain == nil {
		return []entity.UserReferenceDepartment{}, nil
	}
	references, err := whereChain.Preload("User", nil).
		Preload("ReferableDepartment", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ユーザ参照（部署）情報の取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceDepartmentRepository) GetByUserID(userID uint) ([]entity.UserReferenceDepartment, error) {
	references, err := repo.preloadWhere("user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザID指定でのユーザ参照（部署）取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceDepartmentRepository) GetByReferableDepartmentID(referableDepartmentID uint) ([]entity.UserReferenceDepartment, error) {
	references, err := repo.preloadWhere("referable_department_id = ?", referableDepartmentID)
	if err != nil {
		return nil, fmt.Errorf("参照可能部署ID指定でのユーザ参照（部署）取得に失敗しました: %w", err)
	}
	return references, nil
}

func (repo *UserReferenceDepartmentRepository) DeleteByUserID(userID uint) error {
	if _, err := repo.baseRepository.Delete(0, "user_id = ?", userID); err != nil {
		return fmt.Errorf("ユーザ参照（部署）のユーザ単位削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserReferenceDepartmentRepository) IsReferable(userID, referableDepartmentID uint) (isReferable bool, err error) {
	references, err := repo.preloadWhere("user_id = ? AND referable_department_id = ?", userID, referableDepartmentID)
	if err != nil {
		return false, fmt.Errorf("ユーザーIDと参照可能なチームID指定での参照権限有無のチェックに失敗しました: %w", err)
	}
	if len(references) > 0 {
		return true, nil
	}
	return false, nil
}
