package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Name string    `gorm:"type:text;not null;uniqueIndx"`

	DepartmentMember []DepartmentMember `gorm:"foreignKey:RoleID;references:ID"`
}
