package repository_interface

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import (
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
	"time"

	"github.com/google/uuid"
)

// MeetingListFilter is applied in SQL when listing planning meetings.
type MeetingListFilter struct {
	MentorUUIDs       []string
	MenteeUUIDs       []string
	DepartmentUUIDs   []string
	Status            *valueobject.MeetingStatus
	From              *time.Time
	To                *time.Time
	IncludeDeleted    bool
	LimitToMeetingIDs bool
	MeetingIDs        []uint
}

type MeetingRepository interface {
	Create(meeting *entity.Meeting) (*entity.Meeting, error)

	GetAll() ([]entity.Meeting, error)
	GetByID(id uint) (*entity.Meeting, error)
	GetByIDs(id []uint) ([]entity.Meeting, error)
	GetByUUID(uuid uuid.UUID) (*entity.Meeting, error)
	GetByUUIDs(uuids []uuid.UUID) ([]entity.Meeting, error)
	GetByPairingID(pairingID uint) ([]entity.Meeting, error)
	GetByPairingIDs(pairingIDs []uint) ([]entity.Meeting, error)
	GetByPairingIDsByMonth(pairingIDs []uint, month time.Time) ([]entity.Meeting, error)
	GetByFilter(filter MeetingListFilter) ([]entity.Meeting, error)

	Update(meetingID uint, meeting *entity.Meeting) error
	UpdateStatus(meetingID uint, status valueobject.MeetingStatus) error
	UpdateTemplate(meetingID, templateID uint) error
	UpdatePairing(meetingID, pairingID uint) error
	UpdateImplementDate(meetingID uint, implementDate time.Time) error

	GetUntilToday(base time.Time) ([]entity.Meeting, error)
	GetFromSpecificDateOnward(base time.Time) ([]entity.Meeting, error)
	GetBetweenDates(from, to time.Time) ([]entity.Meeting, error)

	Delete(meetingIDs ...uint) error

	GetMeetingIDsBySharerUserID(userID uint) ([]uint, error)
	ReplaceSharers(meetingID uint, userIDs []uint) error
}
