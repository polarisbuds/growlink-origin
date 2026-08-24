package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRelationshipRepository struct {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}

func NewUserRelationshipRepository(ctx context.Context, dp *gorm.DB) *UserRelationshipRepository {
	return &UserRelationshipRepository{
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
