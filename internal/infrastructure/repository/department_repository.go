package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Department]
	unScopedBaseRepository BaseUnScopedRepository[entity.Department]
}

func NewDepartementRepository(ctx context.Context, dp *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{
		baseRepository: BaseRepository[entity.Department]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Department]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}

}
