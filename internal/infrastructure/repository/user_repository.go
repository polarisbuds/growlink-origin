package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.User]
	unScopedBaseRepository BaseUnScopedRepository[entity.User]
}

func NewUserRepository(ctx context.Context, dp *gorm.DB) *UserRepository {
	return &UserRepository{
		baseRepository: BaseRepository[entity.User]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.User]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
