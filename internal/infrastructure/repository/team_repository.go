package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type TeamRepository struct {
	BaseRepository[entity.Team]
	BaseUnScopedRepository[entity.Team]
}

func NewTeamsRepository(ctx context.Context, dp *gorm.DB) *TeamRepository {
	return &TeamRepository{
		BaseRepository: BaseRepository[entity.Team]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.Team]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
