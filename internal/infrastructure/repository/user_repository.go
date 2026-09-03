package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"
	"log/slog"

	"gorm.io/gorm"
)

type UserRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.User]
	unScopedBaseRepository BaseUnScopedRepository[entity.User]
}

func NewUserRepository(ctx context.Context, dp *gorm.DB) *UserRepository {
	return &UserRepository{
		baseRepository: BaseRepository[entity.User]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.User]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *UserRepository) GetByID(id uint) (*entity.User, error) {
	user, err := repo.baseRepository.First("id ASC", "id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("指定されたユーザIDは存在しませんでした: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("指定されたユーザIDは存在しませんでした: %w", gorm.ErrRecordNotFound)
	}
	return user, nil
}

func (repo *UserRepository) GetByIDs(ids []uint) ([]entity.User, error) {
	users, err := repo.baseRepository.WhereInBatch(len(ids), "id ASC", "id IN ?", ids)
	if err != nil {
		return nil, fmt.Errorf("ユーザIDのバッチ取得に失敗しました: %w", err)
	}
	return users, nil
}

func (repo *UserRepository) GetByDepartmentID(departmentID uint) ([]entity.User, error) {
	userWhereChain, err := repo.baseRepository.WhereChain("id ASC", "department_id = ?", departmentID)
	if err != nil {
		return nil, fmt.Errorf("部署ID指定での情報取得用チェインの作成に失敗しました: %w", err)
	}
	if userWhereChain == nil {
		return []entity.User{}, nil
	}
	users, err := userWhereChain.Preload("Department", nil).
		Preload("Role", nil).
		Preload("TeamMembers", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署ID指定でのユーザ情報取得処理に失敗しました: %w", err)
	}
	return users, nil
}

func (repo *UserRepository) UpdateByID(user *entity.User) error {
	if user == nil {
		slog.Warn("アップデート対象のデータが空データ（nil）です。")
		return nil
	}

	if _, err := repo.baseRepository.Updates(user, "id = ?", user.ID); err != nil {
		return fmt.Errorf("ユーザ情報の更新に失敗しました: %w", err)
	}
	return nil
}
