package repository

import (
	"context"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type PeeringRepository struct {
type PairingRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Peering]
	unScopedBaseRepository BaseUnScopedRepository[entity.Peering]
}

func NewPeeringRepository(ctx context.Context, dp *gorm.DB) *PeeringRepository {
	return &PeeringRepository{
		baseRepository: BaseRepository[entity.Peering]{
func NewPairingRepository(ctx context.Context, dp *gorm.DB) *PairingRepository {
	return &PairingRepository{
		baseRepository: BaseRepository[entity.Pairing]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Peering]{
		unScopedBaseRepository: BaseUnScopedRepository[entity.Pairing]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}
