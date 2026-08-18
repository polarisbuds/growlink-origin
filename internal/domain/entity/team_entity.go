package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Name string    `gorm:"type:varchar(64);not null"`

	DepartmentID int64 `gorm:"not null;index"`
	ReaderID     int64 `gorm:"not null;index"`

	Department *Department `gorm:"foreignKey:DepartmentID"`
	Reader     *User       `gorm:"foreignKey:ReaderID"`

	Users []User `gorm:"foreignKey:TeamID"`
}
