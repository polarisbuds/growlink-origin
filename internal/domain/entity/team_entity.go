package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Name string    `gorm:"type:text;not null"`

	DepartmentID int64 `gorm:"not null;index"`

	Department *Department `gorm:"foreignKey:DepartmentID"`

	TeamMembers []TeamMember `gorm:"foreignKey:TeamID"`
}
