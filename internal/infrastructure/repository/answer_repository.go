package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type AnswerRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Answer]
	unScopedBaseRepository BaseUnScopedRepository[entity.Answer]
}

func NewAnswerRepository(ctx context.Context, dp *gorm.DB) *AnswerRepository {
	return &AnswerRepository{
		baseRepository: BaseRepository[entity.Answer]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Answer]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
