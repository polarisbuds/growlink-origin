package repository_interface

import "terms-checker/internal/domain/entity"

type UserRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
