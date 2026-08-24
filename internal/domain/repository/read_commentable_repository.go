package repository_interface

import "growlink/internal/domain/entity"

type ReadCommentableRepository interface {
	BaseRepository[entity.ReadCommentable]
	BaseUnScopedRepository[entity.ReadCommentable]
}
