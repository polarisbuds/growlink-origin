package repository

import (
	"context"
	"terms-checker/internal/domain/entity"

	"gorm.io/gorm"
)

type TopicRepository struct {
	BaseRepository[entity.Topic]
	BaseUnScopedRepository[entity.Topic]
}

func NewTopicRepository(ctx context.Context, dp *gorm.DB) *TopicRepository {
	return &TopicRepository{
		BaseRepository: BaseRepository[entity.Topic]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.Topic]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
