package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type EventRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Event]
	unScopedBaseRepository BaseUnScopedRepository[entity.Event]
}

func NewEventRepository(ctx context.Context, dp *gorm.DB) *EventRepository {
	return &EventRepository{
		baseRepository: BaseRepository[entity.Event]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Event]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
