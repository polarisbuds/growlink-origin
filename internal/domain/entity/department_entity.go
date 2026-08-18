package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Name string    `gorm:"type:varchar(64);not null"`

	Users []User `gorm:"foreignKey:DepartmentID;references:ID"`
	Teams []Team `gorm:"foreignKey:DepartmentID;references:ID"`
}
