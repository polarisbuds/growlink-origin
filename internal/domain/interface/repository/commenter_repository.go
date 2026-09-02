package repository_interface

import "growlink/internal/domain/entity"

type Commenter[T entity.AnswerComment | entity.MeetingComment] interface {
	GetComment(commentID uint) (*T, error)
	PostComment(*T) error
	UpdateComment(commentID uint, newContent string) error
	DeleteComment(commentID uint) error
}
