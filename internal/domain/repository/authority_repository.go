package repository_interface

import "growlink/internal/domain/entity"

type AuthorityRepository interface {
	BaseRepository[entity.Authority]
	BaseUnScopedRepository[entity.Authority]
}
