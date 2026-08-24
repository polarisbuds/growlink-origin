package repository_interface

import "terms-checker/internal/domain/entity"

type MeetingRepository interface {
	BaseRepository[entity.Meeting]
	BaseUnScopedRepository[entity.Meeting]
}
