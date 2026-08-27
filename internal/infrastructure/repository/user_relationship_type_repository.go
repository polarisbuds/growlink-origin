package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRelationshipTypeRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.UserRelationshipType]
	unScopedBaseRepository BaseUnScopedRepository[entity.UserRelationshipType]
}

func NewUserRelationshipTypeRepository(ctx context.Context, dp *gorm.DB) *UserRelationshipTypeRepository {
	return &UserRelationshipTypeRepository{
		baseRepository: BaseRepository[entity.UserRelationshipType]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.UserRelationshipType]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
