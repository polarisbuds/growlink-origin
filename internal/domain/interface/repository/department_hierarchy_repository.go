package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type DepartmentHierarchyRepository interface {
	Create(parentDepartmentID, childDepartmentID uint) error
	Delete(id uint) error
	GetAll() ([]entity.DepartmentHierarchy, error)
	GetByParentDepartmentID(parentDepartmentID uint) ([]entity.DepartmentHierarchy, error)
	GetByChildDepartmentID(childDepartmentID uint) ([]entity.DepartmentHierarchy, error)
	DeleteByParentDepartmentID(parentDepartmentID uint) error
}
