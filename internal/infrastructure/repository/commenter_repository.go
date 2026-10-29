package repository

import (
	"context"
	"errors"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type Commenter[T entity.AnswerComment | entity.MeetingComment] struct {
	ctx context.Context
	db  *gorm.DB
}

func (c *Commenter[T]) GetComment(commentID uint) (*T, error) {
	comment, err := gorm.G[T](c.db).Where("id = ?", commentID).Order("id ASC").First(c.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("コメントの取得に失敗しました: %w", err)
	}
	return &comment, nil
}

func (c *Commenter[T]) PostComment(comment *T) error {
	if err := gorm.G[T](c.db).Create(c.ctx, comment); err != nil {
		return fmt.Errorf("コメントの作成に失敗しました: %w", err)
	}
	return nil
}

func (c *Commenter[T]) UpdateComment(commentID uint, newContent string) error {
	if _, err := gorm.G[T](c.db).Where("id = ?", commentID).Update(c.ctx, "content", newContent); err != nil {
		return fmt.Errorf("コメントの更新に失敗しました: %w", err)
	}
	return nil
}

func (c *Commenter[T]) DeleteComment(commentID uint) error {
	if _, err := gorm.G[T](c.db).Where("id = ?", commentID).Delete(c.ctx); err != nil {
		return fmt.Errorf("コメントの削除に失敗しました: %w", err)
	}
	return nil
}
