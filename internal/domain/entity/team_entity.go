package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Name string    `gorm:"type:text;not null"`

	DepartmentID uint `gorm:"not null;index"`

	Department *Department `gorm:"foreignKey:DepartmentID"`

	TeamMembers []TeamMember `gorm:"foreignKey:TeamID"`

	UserReferenceTeams []UserReferenceTeam `gorm:"foreignKey:ReferableTeamID;references:ID"`
}
