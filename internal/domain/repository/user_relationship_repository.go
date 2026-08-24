package repository_interface

import "terms-checker/internal/domain/entity"

type UserRelationshipRepository interface {
	BaseRepository[entity.UserRelationship]
	BaseUnScopedRepository[entity.UserRelationship]
}
