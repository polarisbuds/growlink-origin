package repository_interface

import "growlink/internal/domain/entity"

type MeetingStatusRepository interface {
	BaseRepository[entity.MeetingStatus]
	BaseUnScopedRepository[entity.MeetingStatus]
}
