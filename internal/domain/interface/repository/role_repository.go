package repository_interface

import (
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
)

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type RoleRepository interface {
	GetAll() ([]entity.Role, error)
	GetByUUID(roleUUID uuid.UUID) (*entity.Role, error)
	GetByName(name string) (*entity.Role, error)
	Create(role *entity.Role) error
	Update(role *entity.Role) error
	Delete(id uint) error
	CountMembersByRoleID(roleID uint) (int64, error)
}
