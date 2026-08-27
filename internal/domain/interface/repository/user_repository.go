package repository_interface

import "growlink/internal/domain/entity"

type UserRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
