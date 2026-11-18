package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type DepartmentDepartmentRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewUserDepartmentRepository(ctx context.Context, dp *gorm.DB) *DepartmentDepartmentRepository {
	return &DepartmentDepartmentRepository{ctx: ctx, db: dp}
}

func (repo *DepartmentDepartmentRepository) GetUserIDsByDepartmentID(departmentID uint) ([]uint, error) {
	var rows []entity.DepartmentMember
	if err := repo.db.WithContext(repo.ctx).Where("department_id = ?", departmentID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("部署所属ユーザの取得に失敗しました: %w", err)
	}
	ids := make([]uint, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.UserID)
	}
	return ids, nil
}

func (repo *DepartmentDepartmentRepository) Create(userID, departmentID uint) error {
	if err := repo.db.WithContext(repo.ctx).Create(&entity.DepartmentMember{
		UserID:       userID,
		DepartmentID: departmentID,
	}).Error; err != nil {
		return fmt.Errorf("部署所属の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentDepartmentRepository) DeleteByUserID(userID uint) error {
	if err := repo.db.WithContext(repo.ctx).Unscoped().Where("user_id = ?", userID).Delete(&entity.DepartmentMember{}).Error; err != nil {
		return fmt.Errorf("部署所属の削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentDepartmentRepository) DeleteByUserAndDepartment(userID, departmentID uint) error {
	if err := repo.db.WithContext(repo.ctx).Unscoped().
		Where("user_id = ? AND department_id = ?", userID, departmentID).
		Delete(&entity.DepartmentMember{}).Error; err != nil {
		return fmt.Errorf("部署所属の削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *DepartmentDepartmentRepository) GetDepartmentIDsByUserID(userID uint) ([]uint, error) {
	var rows []entity.DepartmentMember
	if err := repo.db.WithContext(repo.ctx).Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("ユーザ所属部署の取得に失敗しました: %w", err)
	}
	ids := make([]uint, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.DepartmentID)
	}
	return ids, nil
}

func (repo *DepartmentDepartmentRepository) SyncDepartmentMembers(departmentID uint, userIDs []uint) error {
	want := map[uint]bool{}
	for _, id := range userIDs {
		want[id] = true
	}
	var existing []entity.DepartmentMember
	if err := repo.db.WithContext(repo.ctx).Where("department_id = ?", departmentID).Find(&existing).Error; err != nil {
		return err
	}
	have := map[uint]uint{}
	for _, row := range existing {
		have[row.UserID] = row.ID
	}
	for uid := range have {
		if !want[uid] {
			if err := repo.db.WithContext(repo.ctx).Delete(&entity.DepartmentMember{}, have[uid]).Error; err != nil {
				return err
			}
		}
	}
	for uid := range want {
		if have[uid] == 0 {
			if err := repo.db.WithContext(repo.ctx).Create(&entity.DepartmentMember{
				UserID:       uid,
				DepartmentID: departmentID,
			}).Error; err != nil {
				return err
			}
		}
	}
	for uid := range want {
		user, err := gorm.G[entity.DepartmentMember](repo.db).Where("user_id = ?", uid).First(repo.ctx)
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
