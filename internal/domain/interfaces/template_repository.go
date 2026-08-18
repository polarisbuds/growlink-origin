package interfaces

import "terms-checker/internal/domain/entity"

type TemplateRepository interface {
	BaseRepository[entity.User]
	BaseUnScopedRepository[entity.User]
}
