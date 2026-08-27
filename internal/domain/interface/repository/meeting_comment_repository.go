package repository_interface

import "growlink/internal/domain/entity"

type MeetingCommentRepository interface {
	BaseRepository[entity.MeetingComment]
	BaseUnScopedRepository[entity.MeetingComment]
}
