package interfaces

import "terms-checker/internal/domain/entity"

type TopicRepository interface {
	BaseRepository[entity.Topic]
	BaseUnScopedRepository[entity.Topic]
}
