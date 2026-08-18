package repository

import (
	"context"
	"terms-checker/internal/domain/entity"

	"gorm.io/gorm"
)

type ReadCommentableRepository struct {
	BaseRepository[entity.ReadCommentable]
	BaseUnScopedRepository[entity.ReadCommentable]
}

func NewReadCommentableRepository(ctx context.Context, dp *gorm.DB) *ReadCommentableRepository {
	return &ReadCommentableRepository{
		BaseRepository: BaseRepository[entity.ReadCommentable]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.ReadCommentable]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
