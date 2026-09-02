package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -destination=mock/mock_answer_comment_repository.go -package=mock growlink/internal/domain/interface/repository AnswerCommentRepository

type AnswerCommentRepository interface {
	Commenter[entity.AnswerComment]
}
