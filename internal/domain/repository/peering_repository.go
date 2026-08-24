package repository_interface

import "growlink/internal/domain/entity"

type PeeringRepository interface {
	BaseRepository[entity.Peering]
	BaseUnScopedRepository[entity.Peering]
}
