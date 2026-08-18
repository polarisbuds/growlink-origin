package interfaces

import "terms-checker/internal/domain/entity"

type MeetingCommentRepository interface {
	BaseRepository[entity.MeetingComment]
	BaseUnScopedRepository[entity.MeetingComment]
}
