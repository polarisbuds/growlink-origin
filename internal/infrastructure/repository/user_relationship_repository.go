package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRelationshipRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.User]
	unScopedBaseRepository BaseUnScopedRepository[entity.User]
}

func NewUserRelationshipRepository(ctx context.Context, dp *gorm.DB) *UserRelationshipRepository {
	return &UserRelationshipRepository{
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
