package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type MeetingRepository struct {
	BaseRepository[entity.Meeting]
	BaseUnScopedRepository[entity.Meeting]
}

func NewMeetingRepository(ctx context.Context, dp *gorm.DB) *MeetingRepository {
	return &MeetingRepository{
		BaseRepository: BaseRepository[entity.Meeting]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.Meeting]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
