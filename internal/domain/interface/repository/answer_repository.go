package repository_interface

import "growlink/internal/domain/entity"

type AnswerRepository interface {
	BaseRepository[entity.Answer]
	BaseUnScopedRepository[entity.Answer]
}
