package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -destination=mock/mock_meeting_comment_repository.go -package=mock growlink/internal/domain/interface/repository MeetingCommentRepository

type MeetingCommentRepository interface {
	Commenter[entity.MeetingComment]
}
