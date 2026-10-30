package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type UserReferenceTeamRepository interface {
	Create(userID, referableTeamID uint) error
	Delete(id uint) error
	GetAll() ([]entity.UserReferenceTeam, error)
	GetByUserID(userID uint) ([]entity.UserReferenceTeam, error)
	GetByReferableTeamID(referableTeamID uint) ([]entity.UserReferenceTeam, error)
	IsReferable(userID, targetUserID uint) (isReferable bool, err error)
	DeleteByUserID(userID uint) error
}
