package repository_interface

import "growlink/internal/domain/entity"

type UserRelationshipTypeRepository interface {
	BaseRepository[entity.UserRelationshipType]
	BaseUnScopedRepository[entity.UserRelationshipType]
}
