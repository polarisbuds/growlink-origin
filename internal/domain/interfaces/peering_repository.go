package interfaces

import "terms-checker/internal/domain/entity"

type PeeringRepository interface {
	BaseRepository[entity.Peering]
	BaseUnScopedRepository[entity.Peering]
}
