package repository

import (
	"context"
	"fmt"
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

func (repo *TopicRepository) Create(topic *entity.Topic) error {
	if err := repo.baseRepository.Insert(topic); err != nil {
		return fmt.Errorf("トピックの作成に失敗しました: %w", err)
	}
	return nil
}
