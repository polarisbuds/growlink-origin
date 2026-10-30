package entity

import "gorm.io/gorm"

type UserDepartment struct {
	gorm.Model

	UserID       uint `gorm:"not null;index;uniqueIndex:idx_user_department"`
	DepartmentID uint `gorm:"not null;index;uniqueIndex:idx_user_department"`

	User       *User       `gorm:"foreignKey:UserID"`
	Department *Department `gorm:"foreignKey:DepartmentID"`
}
