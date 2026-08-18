package interfaces

import "terms-checker/internal/domain/entity"

type TeamRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
