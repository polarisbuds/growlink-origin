package interfaces

import "terms-checker/internal/domain/entity"

type AuthorityRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
