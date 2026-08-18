package repository

import (
	"context"
	"terms-checker/internal/domain/entity"

	"gorm.io/gorm"
)

type TemplateRepository struct {
	BaseRepository[entity.Template]
	BaseUnScopedRepository[entity.Template]
}

func NewTemplateRepository(ctx context.Context, dp *gorm.DB) *TemplateRepository {
	return &TemplateRepository{
		BaseRepository: BaseRepository[entity.Template]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.Template]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
