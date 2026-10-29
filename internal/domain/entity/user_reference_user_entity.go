package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserReferenceUser struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	UserID           uint `gorm:"not null;index"`
	ReferableUserID  uint `gorm:"not null;index"`

	User          *User `gorm:"foreignKey:UserID;references:ID"`
	ReferableUser *User `gorm:"foreignKey:ReferableUserID;references:ID"`
}
