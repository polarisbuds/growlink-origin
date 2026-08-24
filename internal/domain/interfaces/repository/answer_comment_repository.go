package repository_interface

import "growlink/internal/domain/entity"

type AnswerCommentRepository interface {
	BaseRepository[entity.AnswerComment]
	BaseUnScopedRepository[entity.AnswerComment]
}
