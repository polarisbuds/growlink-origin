package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type TopicRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Topic]
	unScopedBaseRepository BaseUnScopedRepository[entity.Topic]
}

func NewTopicRepository(ctx context.Context, dp *gorm.DB) *TopicRepository {
	return &TopicRepository{
		baseRepository: BaseRepository[entity.Topic]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Topic]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
