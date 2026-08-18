package repository

import (
	"context"
	"terms-checker/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}

func NewUserRepository(ctx context.Context, dp *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: BaseRepository[entity.User]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.User]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
