package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type UserDepartmentRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewUserDepartmentRepository(ctx context.Context, dp *gorm.DB) *UserDepartmentRepository {
	return &UserDepartmentRepository{ctx: ctx, db: dp}
}

func (repo *UserDepartmentRepository) GetUserIDsByDepartmentID(departmentID uint) ([]uint, error) {
	var rows []entity.UserDepartment
	if err := repo.db.WithContext(repo.ctx).Where("department_id = ?", departmentID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("部署所属ユーザの取得に失敗しました: %w", err)
	}
	ids := make([]uint, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.UserID)
	}
	return ids, nil
}

func (repo *UserDepartmentRepository) GetDepartmentIDsByUserID(userID uint) ([]uint, error) {
	var rows []entity.UserDepartment
	if err := repo.db.WithContext(repo.ctx).Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("ユーザ所属部署の取得に失敗しました: %w", err)
	}
	ids := make([]uint, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.DepartmentID)
	}
	return ids, nil
}

func (repo *UserDepartmentRepository) SyncDepartmentMembers(departmentID uint, userIDs []uint) error {
	want := map[uint]bool{}
	for _, id := range userIDs {
		want[id] = true
	}
	var existing []entity.UserDepartment
	if err := repo.db.WithContext(repo.ctx).Where("department_id = ?", departmentID).Find(&existing).Error; err != nil {
		return err
	}
	have := map[uint]uint{}
	for _, row := range existing {
		have[row.UserID] = row.ID
	}
	for uid := range have {
		if !want[uid] {
			if err := repo.db.WithContext(repo.ctx).Delete(&entity.UserDepartment{}, have[uid]).Error; err != nil {
				return err
			}
		}
	}
	for uid := range want {
		if have[uid] == 0 {
			if err := repo.db.WithContext(repo.ctx).Create(&entity.UserDepartment{
				UserID:       uid,
				DepartmentID: departmentID,
			}).Error; err != nil {
				return err
			}
		}
	}
	for uid := range want {
		user, err := gorm.G[entity.User](repo.db).Where("id = ?", uid).First(repo.ctx)
		if err != nil {
			continue
		}
		if user.DepartmentID == 0 {
			user.DepartmentID = departmentID
			if err := repo.db.WithContext(repo.ctx).Model(&entity.User{}).Where("id = ?", uid).
				Update("department_id", departmentID).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
