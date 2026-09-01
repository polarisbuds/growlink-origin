package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Meeting struct {
	gorm.Model

	UUID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	ImplementedAt time.Time `gorm:"not null"`

	StatusID   uint `gorm:"not null;index"`
	TemplateID uint `gorm:"not null;index"`
	PairingID  uint `gorm:"not null;index"`

	Status   *MeetingStatus `gorm:"foreignKey:StatusID"`
	Template *Template      `gorm:"foreignKey:TemplateID"`
	Pairing  *Pairing       `gorm:"foreignKey:PairingID"`

	Answers          []Answer          `gorm:"foreignKey:MeetingID"`
	MeetingComments  []MeetingComment  `gorm:"foreignKey:MeetingID"`
	ReadCommentables []ReadCommentable `gorm:"foreignKey:MeetingID"`
}
