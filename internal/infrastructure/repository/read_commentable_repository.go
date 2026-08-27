package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type ReadCommentableRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.ReadCommentable]
	unScopedBaseRepository BaseUnScopedRepository[entity.ReadCommentable]
}

func NewReadCommentableRepository(ctx context.Context, dp *gorm.DB) *ReadCommentableRepository {
	return &ReadCommentableRepository{
		baseRepository: BaseRepository[entity.ReadCommentable]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.ReadCommentable]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
