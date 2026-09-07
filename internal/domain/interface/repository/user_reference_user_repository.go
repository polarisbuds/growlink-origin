package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type UserReferenceUserRepository interface {
	Create(userID, referableUserID uint) error
	Delete(id uint) error
	GetAll() ([]entity.UserReferenceUser, error)
	GetByUserID(userID uint) ([]entity.UserReferenceUser, error)
	GetByReferableUserID(referableUserID uint) ([]entity.UserReferenceUser, error)
}
