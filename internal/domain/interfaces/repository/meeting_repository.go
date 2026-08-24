package repository_interface

import "growlink/internal/domain/entity"

type MeetingRepository interface {
	BaseRepository[entity.Meeting]
	BaseUnScopedRepository[entity.Meeting]
}
