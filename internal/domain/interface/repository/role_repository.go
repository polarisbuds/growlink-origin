package repository_interface

import "growlink/internal/domain/entity"

type RoleRepository interface {
	BaseRepository[entity.Role]
	BaseUnScopedRepository[entity.Role]
}
