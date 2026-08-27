package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type AuthorityRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Authority]
	unScopedBaseRepository BaseUnScopedRepository[entity.Authority]
}

func NewAuthorityRepository(ctx context.Context, dp *gorm.DB) *AuthorityRepository {
	return &AuthorityRepository{
		baseRepository: BaseRepository[entity.Authority]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Authority]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
