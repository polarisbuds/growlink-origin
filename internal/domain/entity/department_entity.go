package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Name string    `gorm:"type:text;not null"`

	HeadUserID uint `gorm:"not null;default:0;index"`

	Users []User `gorm:"foreignKey:DepartmentID;references:ID"`
	Teams []Team `gorm:"foreignKey:DepartmentID;references:ID"`

	ParentHierarchies            []DepartmentHierarchy   `gorm:"foreignKey:ParentDepartmentID;references:ID"`
	ChildHierarchies             []DepartmentHierarchy   `gorm:"foreignKey:ChildDepartmentID;references:ID"`
	DepartmentRoles              []DepartmentRole        `gorm:"foreignKey:DepartmentID;references:ID"`
	UserReferenceDepartments     []UserReferenceDepartment `gorm:"foreignKey:ReferableDepartmentID;references:ID"`
}
