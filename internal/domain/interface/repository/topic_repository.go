package repository_interface

import "growlink/internal/domain/entity"

type TopicRepository interface {
	BaseRepository[entity.Topic]
	BaseUnScopedRepository[entity.Topic]
}
