package interfaces

import "terms-checker/internal/domain/entity"

type MeetingStatusRepository interface {
	BaseRepository[entity.MeetingStatus]
	BaseUnScopedRepository[entity.MeetingStatus]
}
