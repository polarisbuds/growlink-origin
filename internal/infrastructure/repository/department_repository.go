package repository

import (
	"context"
	"terms-checker/internal/domain/entity"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	BaseRepository[entity.Department]
	BaseUnScopedRepository[entity.Department]
}

func NewDepartementRepository(ctx context.Context, dp *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{
		BaseRepository: BaseRepository[entity.Department]{
			dp:  dp,
			ctx: ctx,
		},
		BaseUnScopedRepository: BaseUnScopedRepository[entity.Department]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}

}
