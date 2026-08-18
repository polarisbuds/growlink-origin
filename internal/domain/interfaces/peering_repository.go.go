package interfaces

import "terms-checker/internal/domain/entity"

type PeeringRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
