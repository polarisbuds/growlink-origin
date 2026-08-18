package interfaces

import "terms-checker/internal/domain/entity"

type TeamRepository interface {
	BaseRepository[entity.Team]
	BaseUnScopedRepository[entity.Team]
}
