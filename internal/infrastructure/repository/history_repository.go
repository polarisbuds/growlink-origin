package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type HistoryRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.History]
	unScopedBaseRepository BaseUnScopedRepository[entity.History]
}

func NewHistoryRepository(ctx context.Context, dp *gorm.DB) *HistoryRepository {
	return &HistoryRepository{
		baseRepository: BaseRepository[entity.History]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.History]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
