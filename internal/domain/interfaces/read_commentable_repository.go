package interfaces

import "terms-checker/internal/domain/entity"

type ReadCommentableRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
