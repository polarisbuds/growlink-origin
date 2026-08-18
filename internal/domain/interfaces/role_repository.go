package interfaces

import "terms-checker/internal/domain/entity"

type RoleRepository interface {
	BaseRepository[entity.Role]
	BaseUnScopedRepository[entity.Role]
}
