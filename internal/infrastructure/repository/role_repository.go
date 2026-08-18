package repository

import (
	"context"
	"terms-checker/internal/domain/entity"

	"gorm.io/gorm"
)

type RoleRepository struct {
	BaseRepository[entity.Role]
	BaseUnScopedRepository[entity.Role]
}

func NewRoleRepository(ctx context.Context, dp *gorm.DB) *RoleRepository {
	return &RoleRepository{
		BaseRepository: BaseRepository[entity.Role]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.Role]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
