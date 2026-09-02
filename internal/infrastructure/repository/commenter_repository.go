package repository

import (
	"fmt"
	"growlink/internal/domain/entity"
)

type Commenter[T entity.AnswerComment | entity.MeetingComment] struct {
	base BaseRepository[T]
}

func (c *Commenter[T]) GetComment(commentID uint) (*T, error) {
	comment, err := c.base.First("id ASC", "id = ?", commentID)
	if err != nil {
		return nil, fmt.Errorf("コメントの取得に失敗しました: %w", err)
	}
	return comment, nil
}

func (c *Commenter[T]) PostComment(comment *T) error {
	if err := c.base.Insert(comment); err != nil {
		return fmt.Errorf("コメントの作成に失敗しました: %w", err)
	}
	return nil
}

func (c *Commenter[T]) UpdateComment(commentID uint, newContent string) error {
	if _, err := c.base.Update("content", newContent, "id = ?", commentID); err != nil {
		return fmt.Errorf("コメントの更新に失敗しました: %w", err)
	}
	return nil
}

func (c *Commenter[T]) DeleteComment(commentID uint) error {
	if _, err := c.base.Delete(1, "id = ?", commentID); err != nil {
		return fmt.Errorf("コメントの削除に失敗しました: %w", err)
	}
	return nil
}
