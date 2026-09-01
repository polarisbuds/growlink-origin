package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRelationship struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	RelationTypeID uint `gorm:"not null;index"`
	SourceUserID   uint `gorm:"not null;index"`
	TargetUserID   uint `gorm:"not null;index"`

	RelationType *UserRelationshipType `gorm:"foreignKey:RelationTypeID;references:ID"`
	SourceUser   *User                 `gorm:"foreignKey:SourceUserID;references:ID"`
	TargetUser   *User                 `gorm:"foreignKey:TargetUserID;references:ID"`
}
