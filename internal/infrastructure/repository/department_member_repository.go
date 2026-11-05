package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type DepartmentMemberRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewDepartmentMemberRepository(ctx context.Context, dp *gorm.DB) *DepartmentMemberRepository {
	return &DepartmentMemberRepository{ctx: ctx, db: dp}
}

func (repo *DepartmentMemberRepository) preloadWhere(query string, args ...any) ([]entity.DepartmentMember, error) {
	members, err := gorm.G[entity.DepartmentMember](repo.db).
		Where(query, args...).
		Order("id ASC").
		Preload("User", nil).
		Preload("Department", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署メンバー情報の取得に失敗しました: %w", err)
	}
	return members, nil
}

func (repo *DepartmentMemberRepository) GetByDepartmentID(departmentID uint) ([]entity.DepartmentMember, error) {
	members, err := repo.preloadWhere("department_id = ?", departmentID)
	if err != nil {
		return nil, fmt.Errorf("部署ID指定での部署メンバー取得に失敗しました: %w", err)
	}
	return members, nil
}

func (repo *DepartmentMemberRepository) GetByUserID(userID uint) ([]entity.DepartmentMember, error) {
	members, err := repo.preloadWhere("user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザID指定での部署メンバー取得に失敗しました: %w", err)
	}
	return members, nil
}

func (repo *DepartmentMemberRepository) ReplaceByDepartment(departmentID uint, members []entity.DepartmentMember) error {
	err := repo.db.WithContext(repo.ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("department_id = ?", departmentID).Delete(&entity.DepartmentMember{}).Error; err != nil {
			return err
		}
		if len(members) == 0 {
			return nil
		}
		rows := make([]entity.DepartmentMember, 0, len(members))
		for _, member := range members {
			if member.UserID == 0 {
				continue
			}
			rows = append(rows, entity.DepartmentMember{
				DepartmentID: departmentID,
				UserID:       member.UserID,
				IsHead:       member.IsHead,
				IsManager:    member.IsManager,
			})
		}
		if len(rows) == 0 {
			return nil
		}
		return tx.Create(&rows).Error
	})
	if err != nil {
		return fmt.Errorf("部署メンバーの入れ替えに失敗しました: %w", err)
	}
	return nil
}
