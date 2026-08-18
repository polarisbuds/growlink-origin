package interfaces

import "terms-checker/internal/domain/entity"

type MeetingRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
