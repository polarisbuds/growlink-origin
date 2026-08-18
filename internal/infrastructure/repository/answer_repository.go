package repository

import (
	"context"
	"terms-checker/internal/domain/entity"

	"gorm.io/gorm"
)

type AnswerRepository struct {
	BaseRepository[entity.Answer]
	BaseUnScopedRepository[entity.Answer]
}

func NewAnswerRepository(ctx context.Context, dp *gorm.DB) *AnswerRepository {
	return &AnswerRepository{
		BaseRepository: BaseRepository[entity.Answer]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.Answer]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
