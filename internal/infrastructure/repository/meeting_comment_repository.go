package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type MeetingCommentRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.MeetingComment]
	unScopedBaseRepository BaseUnScopedRepository[entity.MeetingComment]
}

func NewMeetingCommentRepository(ctx context.Context, dp *gorm.DB) *MeetingCommentRepository {
	return &MeetingCommentRepository{
		baseRepository: BaseRepository[entity.MeetingComment]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.MeetingComment]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
