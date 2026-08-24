package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MeetingStatus struct {
	gorm.Model

	UUID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Code  int64     `gorm:"not null;uniqueIndex"`
	Label string    `gorm:"type:text;not null;uniqueIndex"`

	Meetings []Meeting `gorm:"foreignKey:StatusID"`
}

func (MeetingStatus) TableName() string {
	return "meeting_status"
}
