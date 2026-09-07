package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepartmentHierarchy struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	ParentDepartmentID uint `gorm:"not null;index"`
	ChildDepartmentID  uint `gorm:"not null;index"`

	ParentDepartment *Department `gorm:"foreignKey:ParentDepartmentID;references:ID"`
	ChildDepartment  *Department `gorm:"foreignKey:ChildDepartmentID;references:ID"`
}
