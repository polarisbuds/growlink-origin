package repository_interface

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import (
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(id uint) (*entity.User, error)
	GetByUUID(userUUID uuid.UUID) (*entity.User, error)
	AttachImageBase64(user *entity.User)
	GetByEmail(email string) (*entity.User, error)
	GetByIDs(ids []uint) ([]entity.User, error)
	GetByDepartmentID(departmentID uint) ([]entity.User, error)
	GetAll() ([]entity.User, error)
	SetSuspendedByIDs(ids []uint, suspended bool) error
	Create(user *entity.User) error
	UpdateByID(user *entity.User) error
}
