package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type MeetingCommentRepository struct {
	Commenter[entity.MeetingComment]
}

func NewMeetingCommentRepository(ctx context.Context, dp *gorm.DB) *MeetingCommentRepository {
	return &MeetingCommentRepository{
		Commenter: Commenter[entity.MeetingComment]{
			ctx: ctx,
			db:  dp,
		},
	}
}
