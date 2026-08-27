package repository_interface

import "growlink/internal/domain/entity"

type TeamRepository interface {
	BaseRepository[entity.Team]
	BaseUnScopedRepository[entity.Team]
}
