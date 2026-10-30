package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserReferenceTeam struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	UserID          uint `gorm:"not null;index"`
	ReferableTeamID uint `gorm:"not null;index"`

	User          *User `gorm:"foreignKey:UserID;references:ID"`
	ReferableTeam *Team `gorm:"foreignKey:ReferableTeamID;references:ID"`
}
