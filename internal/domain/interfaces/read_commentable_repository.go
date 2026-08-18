package interfaces

import "terms-checker/internal/domain/entity"

type ReadCommentableRepository interface {
	BaseRepository[entity.ReadCommentable]
	BaseUnScopedRepository[entity.ReadCommentable]
}
