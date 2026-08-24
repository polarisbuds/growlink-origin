package repository_interface

import "terms-checker/internal/domain/entity"

type AnswerCommentRepository interface {
	BaseRepository[entity.AnswerComment]
	BaseUnScopedRepository[entity.AnswerComment]
}
