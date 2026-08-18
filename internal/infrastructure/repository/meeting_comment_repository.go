package repository

import (
	"context"
	"terms-checker/internal/domain/entity"

	"gorm.io/gorm"
)

type MeetingCommentRepository struct {
	BaseRepository[entity.MeetingComment]
	BaseUnScopedRepository[entity.MeetingComment]
}

func NewMeetingCommentRepository(ctx context.Context, dp *gorm.DB) *MeetingCommentRepository {
	return &MeetingCommentRepository{
		BaseRepository: BaseRepository[entity.MeetingComment]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.MeetingComment]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
