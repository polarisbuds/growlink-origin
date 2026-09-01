package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pairing struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	MentorID int64 `gorm:"not null;index"`
	MenteeID int64 `gorm:"not null;index"`

	Mentor *User `gorm:"foreignKey:MentorID"`
	Mentee *User `gorm:"foreignKey:MenteeID"`

	Meetings []Meeting `gorm:"foreignKey:PairingID"`
}
