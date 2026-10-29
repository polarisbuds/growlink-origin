package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type TopicRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewTopicRepository(ctx context.Context, dp *gorm.DB) *TopicRepository {
	return &TopicRepository{ctx: ctx, db: dp}
}

func (repo *TopicRepository) Create(topic *entity.Topic) error {
	if err := gorm.G[entity.Topic](repo.db).Create(repo.ctx, topic); err != nil {
		return fmt.Errorf("トピックの作成に失敗しました: %w", err)
	}
	return nil
}
