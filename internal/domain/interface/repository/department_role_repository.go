package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type DepartmentRoleRepository interface {
	Create(userID, departmentID uint, name string) error
	Update(id uint, role *entity.DepartmentRole) error
	Delete(id uint) error
	GetAll() ([]entity.DepartmentRole, error)
	GetByUserID(userID uint) ([]entity.DepartmentRole, error)
	GetByDepartmentID(departmentID uint) ([]entity.DepartmentRole, error)
}
