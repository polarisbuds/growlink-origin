package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	UUID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Name     string    `gorm:"type:text;not null"`
	Email    string    `gorm:"type:text;not null;uniqueIndex"`
	Password *string   `gorm:"type:text"`

	RoleID       int64 `gorm:"not null;index"`
	DepartmentID int64 `gorm:"not null;index"`
	TeamID       int64 `gorm:"not null;index"`
	AuthorityID  int64 `gorm:"not null;index"`

	Role       *Role       `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Department *Department `gorm:"foreignKey:DepartmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Team       *Team       `gorm:"foreignKey:TeamID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Authority  *Authority  `gorm:"foreignKey:AuthorityID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
