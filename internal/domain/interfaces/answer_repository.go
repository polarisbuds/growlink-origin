package interfaces

import "terms-checker/internal/domain/entity"

type AnswerRepository interface {
	BaseRepository[entity.Answer]
	BaseUnScopedRepository[entity.Answer]
}
