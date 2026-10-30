package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type AnswerCommentRepository struct {
	Commenter[entity.AnswerComment]
}

func NewAnswerCommentRepository(ctx context.Context, dp *gorm.DB) *AnswerCommentRepository {
	return &AnswerCommentRepository{
		Commenter: Commenter[entity.AnswerComment]{
			ctx: ctx,
			db:  dp,
		},
	}
}
