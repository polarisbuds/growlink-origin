package repository_interface

import "growlink/internal/domain/entity"

type UserRepository interface {
	GetByID(id int64) (*entity.User, error)
	UpdateByID(user *entity.User) error
}
