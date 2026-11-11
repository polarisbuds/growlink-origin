package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepartmentMember struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	IsHead    bool `gorm:"not null;default:false"`
	IsManager bool `gorm:"not null;default:false"`

	DepartmentID uint `gorm:"not null;uniqueIndex:idx_department_member"`
	UserID       uint `gorm:"not null;uniqueIndex:idx_department_member"`
	RoleID       uint `gorm:"not null;"`

	Department *Department `gorm:"foreignKey:DepartmentID;references:ID"`
	User       *User       `gorm:"foreignKey:UserID;references:ID"`
	Role       *Role       `gorm:"foreignKey:RoleID;references:ID"`
}
