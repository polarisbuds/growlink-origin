package repository_interface

import "terms-checker/internal/domain/entity"

type AuthorityRepository interface {
	BaseRepository[entity.Authority]
	BaseUnScopedRepository[entity.Authority]
}
