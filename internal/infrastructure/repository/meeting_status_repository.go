package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type MeetingStatusRepository struct {
	BaseRepository[entity.MeetingStatus]
	BaseUnScopedRepository[entity.MeetingStatus]
}

func NewMeetingStatusRepository(ctx context.Context, dp *gorm.DB) *MeetingStatusRepository {
	return &MeetingStatusRepository{
		BaseRepository: BaseRepository[entity.MeetingStatus]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.MeetingStatus]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
