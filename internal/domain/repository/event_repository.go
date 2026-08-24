package repository_interface

import "growlink/internal/domain/entity"

type EventRepository interface {
	BaseRepository[entity.Event]
	BaseUnScopedRepository[entity.Event]
}
