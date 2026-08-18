package interfaces

import "terms-checker/internal/domain/entity"

type AnswerRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
