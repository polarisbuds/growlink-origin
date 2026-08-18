package interfaces

import "terms-checker/internal/domain/entity"

type MeetingCommentRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
