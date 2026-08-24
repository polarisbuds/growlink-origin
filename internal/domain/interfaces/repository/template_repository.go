package repository_interface

import "growlink/internal/domain/entity"

type TemplateRepository interface {
	BaseRepository[entity.Template]
	BaseUnScopedRepository[entity.Template]
}
