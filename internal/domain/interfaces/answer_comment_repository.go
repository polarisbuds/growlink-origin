package interfaces

import "terms-checker/internal/domain/entity"

type AnswerCommentRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
