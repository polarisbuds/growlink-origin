package repository_interface

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import "growlink/internal/domain/entity"

type UserRepository interface {
	GetByID(id uint) (*entity.User, error)
	GetByIDs(ids []uint) ([]entity.User, error)
	UpdateByID(user *entity.User) error
}
