package interfaces

import "terms-checker/internal/domain/entity"

type MeetingStatusRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
