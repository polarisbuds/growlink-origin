package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type MeetingCommentRepository struct {
	ctx context.Context
	Commenter[entity.MeetingComment]
}

func NewMeetingCommentRepository(ctx context.Context, dp *gorm.DB) *MeetingCommentRepository {
	return &MeetingCommentRepository{
		Commenter: Commenter[entity.MeetingComment]{
			BaseRepository[entity.MeetingComment]{
				dp:   dp,
				bctx: ctx,
			},
		},
	}
}
