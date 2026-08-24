package repository_interface

import "growlink/internal/domain/entity"

type UserRelationshipRepository interface {
	BaseRepository[entity.UserRelationship]
	BaseUnScopedRepository[entity.UserRelationship]
}
