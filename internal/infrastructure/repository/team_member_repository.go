package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type TeamMemberRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.TeamMember]
	unScopedBaseRepository BaseUnScopedRepository[entity.TeamMember]
}

func NewTeamMemberRepository(ctx context.Context, dp *gorm.DB) *TeamMemberRepository {
	return &TeamMemberRepository{
		baseRepository: BaseRepository[entity.TeamMember]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.TeamMember]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
