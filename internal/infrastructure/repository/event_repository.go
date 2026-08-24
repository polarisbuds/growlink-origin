package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type EventRepository struct {
	BaseRepository[entity.Event]
	BaseUnScopedRepository[entity.Event]
}

func NewEventRepository(ctx context.Context, dp *gorm.DB) *EventRepository {
	return &EventRepository{
		BaseRepository: BaseRepository[entity.Event]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.Event]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
