package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type HistoryRepository struct {
	BaseRepository[entity.History]
	BaseUnScopedRepository[entity.History]
}

func NewHistoryRepository(ctx context.Context, dp *gorm.DB) *HistoryRepository {
	return &HistoryRepository{
		BaseRepository: BaseRepository[entity.History]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.History]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
