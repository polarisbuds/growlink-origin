package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type TeamMemberRepository struct {
	BaseRepository[entity.TeamMember]
	BaseUnScopedRepository[entity.TeamMember]
}

func NewTeamMemberRepository(ctx context.Context, dp *gorm.DB) *TeamMemberRepository {
	return &TeamMemberRepository{
		BaseRepository: BaseRepository[entity.TeamMember]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.TeamMember]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
