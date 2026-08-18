package repository

import (
	"context"
	"terms-checker/internal/domain/entity"

	"gorm.io/gorm"
)

type AnswerCommentRepository struct {
	BaseRepository[entity.AnswerComment]
	BaseUnScopedRepository[entity.AnswerComment]
}

func NewAnswerCommentRepository(ctx context.Context, dp *gorm.DB) *AnswerCommentRepository {
	return &AnswerCommentRepository{
		BaseRepository: BaseRepository[entity.AnswerComment]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.AnswerComment]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
