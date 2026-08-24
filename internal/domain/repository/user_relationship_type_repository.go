package repository_interface

import "terms-checker/internal/domain/entity"

type UserRelationshipTypeRepository interface {
	BaseRepository[entity.UserRelationshipType]
	BaseUnScopedRepository[entity.UserRelationshipType]
}
