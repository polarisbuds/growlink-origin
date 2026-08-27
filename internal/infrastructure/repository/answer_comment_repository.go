package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type AnswerCommentRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.AnswerComment]
	unScopedBaseRepository BaseUnScopedRepository[entity.AnswerComment]
}

func NewAnswerCommentRepository(ctx context.Context, dp *gorm.DB) *AnswerCommentRepository {
	return &AnswerCommentRepository{
		baseRepository: BaseRepository[entity.AnswerComment]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.AnswerComment]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
