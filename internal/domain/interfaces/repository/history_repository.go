package repository_interface

import "growlink/internal/domain/entity"

type HistoryRepository interface {
	BaseRepository[entity.History]
	BaseUnScopedRepository[entity.History]
}
