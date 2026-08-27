package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type PeeringRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Peering]
	unScopedBaseRepository BaseUnScopedRepository[entity.Peering]
}

func NewPeeringRepository(ctx context.Context, dp *gorm.DB) *PeeringRepository {
	return &PeeringRepository{
		baseRepository: BaseRepository[entity.Peering]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Peering]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
