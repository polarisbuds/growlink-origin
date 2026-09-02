package repository_interface

import (
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
)

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type UserRelationshipRepository interface {
	Get(relationshipID uint) (*entity.UserRelationship, error)
	GetByRelationType(relationshipType valueobject.RelationshipType) ([]entity.UserRelationship, error)
	GetByBossID(bossID uint) (subordinates []entity.UserRelationship, err error)
	GetBySubordinateID(subordinateID uint) (bosses []entity.UserRelationship, err error)
	GetByMenteeID(menteeID uint) (mentors []entity.UserRelationship, err error)
	GetByMentorID(mentorID uint) (mentees []entity.UserRelationship, err error)
	Create(sourceUserID, targetUserID uint, relationshipType valueobject.RelationshipType) error
	Update(relationship *entity.UserRelationship) error
	Delete(relationshipID uint) error
}
