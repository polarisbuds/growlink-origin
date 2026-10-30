package repository

import (
	"context"
	"errors"
	"fmt"
	"growlink/internal/domain/entity"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	ctx       context.Context
	db        *gorm.DB
	uploadDir string
}

func NewUserRepository(ctx context.Context, dp *gorm.DB, uploadDir string) *UserRepository {
	return &UserRepository{ctx: ctx, db: dp, uploadDir: uploadDir}
}

func (repo *UserRepository) AttachImageBase64(user *entity.User) {
	if user == nil || user.ImageURI == nil || *user.ImageURI == "" {
		return
	}
	dataURL, err := readImageDataURL(repo.uploadDir, *user.ImageURI)
	if err != nil {
		slog.Warn("ユーザアイコンのBase64変換に失敗しました", "filename", *user.ImageURI, "err", err)
		return
	}
	user.ImageBase64 = dataURL
}

func (repo *UserRepository) GetByEmail(email string) (*entity.User, error) {
	user, err := gorm.G[entity.User](repo.db).Where("email = ?", email).Order("id ASC").First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) Create(user *entity.User) error {
	if user == nil {
		return fmt.Errorf("作成対象のユーザが空です")
	}
	if err := gorm.G[entity.User](repo.db).Create(repo.ctx, user); err != nil {
		return fmt.Errorf("ユーザの作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserRepository) GetByID(id uint) (*entity.User, error) {
	user, err := gorm.G[entity.User](repo.db).Where("id = ?", id).Order("id ASC").First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("指定されたユーザIDは存在しませんでした: %w", err)
	}
	return &user, nil
}

func (repo *UserRepository) GetByUUID(userUUID uuid.UUID) (*entity.User, error) {
	user, err := gorm.G[entity.User](repo.db).Where("uuid = ?", userUUID.String()).Order("id ASC").First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("指定されたユーザUUIDは存在しませんでした: %w", err)
	}
	repo.AttachImageBase64(&user)
	return &user, nil
}

func (repo *UserRepository) GetByIDs(ids []uint) ([]entity.User, error) {
	if len(ids) == 0 {
		return []entity.User{}, nil
	}
	users, err := gorm.G[entity.User](repo.db).Where("id IN ?", ids).Order("id ASC").Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ユーザIDのバッチ取得に失敗しました: %w", err)
	}
	return users, nil
}

func (repo *UserRepository) GetByDepartmentID(departmentID uint) ([]entity.User, error) {
	var junctionIDs []uint
	if err := repo.db.WithContext(repo.ctx).
		Model(&entity.UserDepartment{}).
		Where("department_id = ?", departmentID).
		Pluck("user_id", &junctionIDs).Error; err != nil {
		return nil, fmt.Errorf("部署所属ユーザの取得に失敗しました: %w", err)
	}
	chain := gorm.G[entity.User](repo.db).Order("id ASC")
	if len(junctionIDs) == 0 {
		chain = chain.Where("department_id = ?", departmentID)
	} else {
		chain = chain.Where("department_id = ? OR id IN ?", departmentID, junctionIDs)
	}
	users, err := chain.
		Preload("Department", nil).
		Preload("Role", nil).
		Preload("TeamMembers", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署ID指定でのユーザ情報取得処理に失敗しました: %w", err)
	}
	return users, nil
}

func (repo *UserRepository) GetAll() ([]entity.User, error) {
	users, err := gorm.G[entity.User](repo.db).Order("id ASC").
		Preload("Department", nil).
		Preload("Role", nil).
		Preload("TeamMembers", nil).
		Preload("TeamMembers.Team", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ユーザの全取得に失敗しました: %w", err)
	}
	return users, nil
}

func (repo *UserRepository) SetSuspendedByIDs(ids []uint, suspended bool) error {
	if len(ids) == 0 {
		return nil
	}
	if _, err := gorm.G[entity.User](repo.db).Where("id IN ?", ids).Update(repo.ctx, "suspended", suspended); err != nil {
		return fmt.Errorf("ユーザの停止状態の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserRepository) UpdateByID(user *entity.User) error {
	if user == nil {
		slog.Warn("アップデート対象のデータが空データ（nil）です。")
		return nil
	}

	if _, err := gorm.G[entity.User](repo.db).Where("id = ?", user.ID).Omit(clause.Associations).Updates(repo.ctx, *user); err != nil {
		return fmt.Errorf("ユーザ情報の更新に失敗しました: %w", err)
	}
	return nil
}
