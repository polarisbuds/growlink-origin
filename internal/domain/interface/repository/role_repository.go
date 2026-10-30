package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type RoleRepository interface {
	GetAll() ([]entity.Role, error)
	Create(role *entity.Role) error
}
