package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserReferenceDepartment struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	UserID                uint `gorm:"not null;index"`
	ReferableDepartmentID uint `gorm:"not null;index"`

	User                *User       `gorm:"foreignKey:UserID;references:ID"`
	ReferableDepartment *Department `gorm:"foreignKey:ReferableDepartmentID;references:ID"`
}
