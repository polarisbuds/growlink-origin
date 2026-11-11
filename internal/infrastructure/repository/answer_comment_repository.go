package repository

import (
	"context"
	"errors"
	"fmt"
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

func (repo *AnswerCommentRepository) GetByAnswerID(answerID uint) (*entity.AnswerComment, error) {
	comment, err := gorm.G[entity.AnswerComment](repo.db).Where("answer_id = ?", answerID).Order("id ASC").First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("回答コメントの取得に失敗しました: %w", err)
	}
	return &comment, nil
}
