package repository_interface

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import (
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
	"time"

	"github.com/google/uuid"
)

type MeetingRepository interface {
	Create(meeting *entity.Meeting) (*entity.Meeting, error)

	GetAll() ([]entity.Meeting, error)
	GetByID(id uint) (*entity.Meeting, error)
	GetByUUID(uuid uuid.UUID) (*entity.Meeting, error)
	GetByPairingID(pairingID uint) ([]entity.Meeting, error)

	Update(meetingID uint, meeting *entity.Meeting) error
	UpdateStatus(meetingID uint, status valueobject.MeetingStatus) error
	UpdateTemplate(meetingID, templateID uint) error
	UpdatePairing(meetingID, pairingID uint) error
	UpdateImplementDate(meetingID uint, implementDate time.Time) error

	GetUntilToday(base time.Time) ([]entity.Meeting, error)
	GetFromSpecificDateOnward(base time.Time) ([]entity.Meeting, error)
	GetBetweenDates(from, to time.Time) ([]entity.Meeting, error)

	Delete(meetingIDs ...uint) error
}
