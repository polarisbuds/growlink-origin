package repository_interface

import "terms-checker/internal/domain/entity"

type TemplateRepository interface {
	BaseRepository[entity.Template]
	BaseUnScopedRepository[entity.Template]
}
