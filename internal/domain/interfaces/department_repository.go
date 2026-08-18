package interfaces

import "terms-checker/internal/domain/entity"

type DepartmentRepository interface {
	BaseRepository[entity.Department]
	BaseUnScopedRepository[entity.Department]
}
