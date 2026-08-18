package repository

import (
	"context"
	"terms-checker/internal/domain/entity"

	"gorm.io/gorm"
)

type AuthorityRepository struct {
	BaseRepository[entity.Authority]
	BaseUnScopedRepository[entity.Authority]
}

func NewAuthorityRepository(ctx context.Context, dp *gorm.DB) *AuthorityRepository {
	return &AuthorityRepository{
		BaseRepository: BaseRepository[entity.Authority]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.Authority]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
