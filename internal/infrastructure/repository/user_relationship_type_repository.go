package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRelationshipTypeRepository struct {
	BaseRepository[entity.UserRelationshipType]
	BaseUnScopedRepository[entity.UserRelationshipType]
}

func NewUserRelationshipTypeRepository(ctx context.Context, dp *gorm.DB) *UserRelationshipTypeRepository {
	return &UserRelationshipTypeRepository{
		BaseRepository: BaseRepository[entity.UserRelationshipType]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.UserRelationshipType]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
