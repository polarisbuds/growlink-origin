package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type MeetingStatusRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.MeetingStatus]
	unScopedBaseRepository BaseUnScopedRepository[entity.MeetingStatus]
}

func NewMeetingStatusRepository(ctx context.Context, dp *gorm.DB) *MeetingStatusRepository {
	return &MeetingStatusRepository{
		baseRepository: BaseRepository[entity.MeetingStatus]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.MeetingStatus]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
