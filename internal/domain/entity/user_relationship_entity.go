package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRelationship struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	RelationTypeID int64 `gorm:"not null;index"`
	SourceUserID   int64 `gorm:"not null;index"`
	TargetUserID   int64 `gorm:"not null;index"`

	RelationType *UserRelationshipType `gorm:"foreignKey:RelationTypeID"`
	SourceUser   *User                 `gorm:"foreignKey:SourceUserID"`
	TargetUser   *User                 `gorm:"foreignKey:TargetUserID"`
}
