package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type DepartmentRepository interface {
	GetByID(id uint) (*entity.Department, error)
	GetAll() ([]entity.Department, error)
	Create(department *entity.Department) error
	UpdateByID(department *entity.Department) error
	Delete(id uint) error
}
