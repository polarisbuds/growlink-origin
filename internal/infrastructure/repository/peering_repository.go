package repository

import (
	"context"
	"terms-checker/internal/domain/entity"

	"gorm.io/gorm"
)

type PeeringRepository struct {
	BaseRepository[entity.Peering]
	BaseUnScopedRepository[entity.Peering]
}

func NewPeeringRepository(ctx context.Context, dp *gorm.DB) *PeeringRepository {
	return &PeeringRepository{
		BaseRepository: BaseRepository[entity.Peering]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.Peering]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
