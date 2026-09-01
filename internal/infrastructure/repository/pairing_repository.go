package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type PairingRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Pairing]
	unScopedBaseRepository BaseUnScopedRepository[entity.Pairing]
}

func NewPairingRepository(ctx context.Context, dp *gorm.DB) *PairingRepository {
	return &PairingRepository{
		baseRepository: BaseRepository[entity.Pairing]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Pairing]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
