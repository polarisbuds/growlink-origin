package repository_interface

//go:generate mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import (
	"growlink/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type MeetingRepository interface {
	GetAll() ([]entity.Meeting, error)
	GetByID(id int64) (*entity.Meeting, error)
	GetByUUID(uuid uuid.UUID) (*entity.Meeting, error)

	GetUntilToday(base time.Time) ([]entity.Meeting, error)
	GetFromSpecificDateOnward(base time.Time) ([]entity.Meeting, error)
	GetBetweenDates(from, to time.Time) ([]entity.Meeting, error)
}
