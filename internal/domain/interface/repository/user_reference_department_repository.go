package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type UserReferenceDepartmentRepository interface {
	Create(userID, referableDepartmentID uint) error
	Delete(id uint) error
	GetAll() ([]entity.UserReferenceDepartment, error)
	GetByUserID(userID uint) ([]entity.UserReferenceDepartment, error)
	GetByReferableDepartmentID(referableDepartmentID uint) ([]entity.UserReferenceDepartment, error)
	IsReferable(userID, referableDepartmentID uint) (isReferable bool, err error)
}
