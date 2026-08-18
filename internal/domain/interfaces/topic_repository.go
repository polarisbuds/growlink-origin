package interfaces

import "terms-checker/internal/domain/entity"

type TopicRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
