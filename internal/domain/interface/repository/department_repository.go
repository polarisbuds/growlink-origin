package repository_interface

import (
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
)

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type DepartmentRepository interface {
	GetByID(id uint) (*entity.Department, error)
	GetByUUID(deptUUID uuid.UUID) (*entity.Department, error)
	GetAll() ([]entity.Department, error)
	Create(department *entity.Department) error
	UpdateByID(department *entity.Department) error
	UpdateHeadUserID(id, headUserID uint) error
	Delete(id uint) error
}
