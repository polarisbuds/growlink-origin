package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepartmentRole struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	Name string `gorm:"type:text;not null"`

	UserID       uint `gorm:"not null;index"`
	DepartmentID uint `gorm:"not null;index"`

	User       *User       `gorm:"foreignKey:UserID;references:ID"`
	Department *Department `gorm:"foreignKey:DepartmentID;references:ID"`
}
