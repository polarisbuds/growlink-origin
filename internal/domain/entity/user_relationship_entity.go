package entity

import (
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRelationship struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	RelationType valueobject.RelationshipType `gorm:"type:smallint;not null;index"`
	SourceUserID uint                         `gorm:"not null;index"`
	TargetUserID uint                         `gorm:"not null;index"`

	SourceUser *User `gorm:"foreignKey:SourceUserID;references:ID"`
	TargetUser *User `gorm:"foreignKey:TargetUserID;references:ID"`
}
