package repository_interface

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import "growlink/internal/domain/entity"

type UserRepository interface {
	GetByID(id int64) (*entity.User, error)
	UpdateByID(user *entity.User) error
}
