package repository_interface

import "growlink/internal/domain/entity"

type DepartmentRepository interface {
	BaseRepository[entity.Department]
	BaseUnScopedRepository[entity.Department]
}
