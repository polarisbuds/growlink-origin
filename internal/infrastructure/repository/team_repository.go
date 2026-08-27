package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type TeamRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Team]
	unScopedBaseRepository BaseUnScopedRepository[entity.Team]
}

func NewTeamsRepository(ctx context.Context, dp *gorm.DB) *TeamRepository {
	return &TeamRepository{
		baseRepository: BaseRepository[entity.Team]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Team]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
