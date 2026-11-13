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
		Preload("Role", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署メンバー情報の取得に失敗しました: %w", err)
	}
	return members, nil
}

func (repo *DepartmentMemberRepository) GetAll() ([]entity.DepartmentMember, error) {
	members, err := gorm.G[entity.DepartmentMember](repo.db).Order("id ASC").Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部署メンバーの全取得に失敗しました: %w", err)
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
	want := map[uint]entity.DepartmentMember{}
	for _, member := range members {
		if member.UserID == 0 {
			continue
		}
		want[member.UserID] = entity.DepartmentMember{
			DepartmentID: departmentID,
			UserID:       member.UserID,
			RoleID:       member.RoleID,
			IsHead:       member.IsHead,
			IsManager:    member.IsManager,
		}
	}

	err := repo.db.WithContext(repo.ctx).Transaction(func(tx *gorm.DB) error {
		var existing []entity.DepartmentMember
		if err := tx.Where("department_id = ?", departmentID).Find(&existing).Error; err != nil {
			return err
		}
		byUser := map[uint]entity.DepartmentMember{}
		for _, row := range existing {
			byUser[row.UserID] = row
		}

		for userID, row := range byUser {
			if _, ok := want[userID]; !ok {
				if err := tx.Unscoped().Delete(&row).Error; err != nil {
					return err
				}
			}
		}
		for userID, desired := range want {
			current, ok := byUser[userID]
			if !ok {
				if err := tx.Create(&desired).Error; err != nil {
					return err
				}
				continue
			}
			if current.RoleID == desired.RoleID && current.IsHead == desired.IsHead && current.IsManager == desired.IsManager {
				continue
			}
			if err := tx.Model(&entity.DepartmentMember{}).
				Where("id = ?", current.ID).
				Select("RoleID", "IsHead", "IsManager").
				Updates(map[string]any{
					"RoleID":    desired.RoleID,
					"IsHead":    desired.IsHead,
					"IsManager": desired.IsManager,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("部署メンバーの入れ替えに失敗しました: %w", err)
	}
	return nil
}
