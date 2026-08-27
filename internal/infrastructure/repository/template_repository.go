package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type TemplateRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Template]
	unScopedBaseRepository BaseUnScopedRepository[entity.Template]
}

func NewTemplateRepository(ctx context.Context, dp *gorm.DB) *TemplateRepository {
	return &TemplateRepository{
		baseRepository: BaseRepository[entity.Template]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Template]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
