package repository_interface

import (
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
)

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type TeamRepository interface {
	GetByID(id uint) (*entity.Team, error)
	GetByUUID(teamUUID uuid.UUID) (*entity.Team, error)
	GetAll() ([]entity.Team, error)
	GetByDepartmentID(departmentID uint) ([]entity.Team, error)
	Create(team *entity.Team) error
	UpdateByID(team *entity.Team) error
	Delete(id uint) error
}
