package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRelationshipRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.UserRelationship]
	unScopedBaseRepository BaseUnScopedRepository[entity.UserRelationship]
}

func NewUserRelationshipRepository(ctx context.Context, dp *gorm.DB) *UserRelationshipRepository {
	return &UserRelationshipRepository{
		baseRepository: BaseRepository[entity.UserRelationship]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.UserRelationship]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
