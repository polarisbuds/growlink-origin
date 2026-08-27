package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type RoleRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Role]
	unScopedBaseRepository BaseUnScopedRepository[entity.Role]
}

func NewRoleRepository(ctx context.Context, dp *gorm.DB) *RoleRepository {
	return &RoleRepository{
		baseRepository: BaseRepository[entity.Role]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Role]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
