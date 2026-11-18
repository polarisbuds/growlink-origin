package repository

import (
	"context"
	"errors"
	"fmt"
	"growlink/internal/domain/entity"
	logutils "growlink/internal/shared/logUtils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewRoleRepository(ctx context.Context, dp *gorm.DB) *RoleRepository {
	return &RoleRepository{ctx: ctx, db: dp}
}

func (repo *RoleRepository) GetByUUID(roleUUID uuid.UUID) (ret0 *entity.Role, err error) {
	const op = "role_repository.GetByUUID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	role, err := gorm.G[entity.Role](repo.db).Where("uuid = ?", roleUUID).First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("UUID指定での役職取得に失敗しました: %w", err)
	}
	return &role, nil
}

func (repo *RoleRepository) GetByUUIDs(roleUUID []string) (ret0 []entity.Role, err error) {
	const op = "role_repository.GetByUUID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	roles, err := gorm.G[entity.Role](repo.db).Where("uuid IN ?", roleUUID).Find(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("UUID指定での役職取得に失敗しました: %w", err)
	}
	return roles, nil
}

func (repo *RoleRepository) GetAll() (ret0 []entity.Role, err error) {
	const op = "role_repository.GetAll"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	roles, err := gorm.G[entity.Role](repo.db).Order("id ASC").Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ロールの全取得に失敗しました: %w", err)
	}
	return roles, nil
}

func (repo *RoleRepository) GetByName(name string) (ret0 *entity.Role, err error) {
	const op = "role_repository.GetByName"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	role, err := gorm.G[entity.Role](repo.db).Where("name = ?", name).First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("役職名での取得に失敗しました: %w", err)
	}
	return &role, nil
}

func (repo *RoleRepository) Create(role *entity.Role) (err error) {
	const op = "role_repository.Create"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if role == nil {
		return fmt.Errorf("役職が空です")
	}
	if err = gorm.G[entity.Role](repo.db).Create(repo.ctx, role); err != nil {
		return fmt.Errorf("役職の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *RoleRepository) Update(role *entity.Role) (err error) {
	const op = "role_repository.Update"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if role == nil || role.ID == 0 {
		return fmt.Errorf("役職が空です")
	}
	if _, err = gorm.G[entity.Role](repo.db).Where("id = ?", role.ID).Updates(repo.ctx, *role); err != nil {
		return fmt.Errorf("役職の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *RoleRepository) Delete(id uint) (err error) {
	const op = "role_repository.Delete"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if _, err = gorm.G[entity.Role](repo.db).Where("id = ?", id).Delete(repo.ctx); err != nil {
		return fmt.Errorf("役職の削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *RoleRepository) CountMembersByRoleID(roleID uint) (ret0 int64, err error) {
	const op = "role_repository.CountMembersByRoleID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	var count int64
	if err = repo.db.WithContext(repo.ctx).Model(&entity.DepartmentMember{}).Where("role_id = ?", roleID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("役職の割当件数取得に失敗しました: %w", err)
	}
	return count, nil
}
