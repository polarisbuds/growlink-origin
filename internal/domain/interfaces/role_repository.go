package interfaces

import "terms-checker/internal/domain/entity"

type RoleRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
