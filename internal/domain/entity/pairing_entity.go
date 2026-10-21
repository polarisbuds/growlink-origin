package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pairing struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	MentorID        uint `gorm:"not null;index"`
	MenteeID        uint `gorm:"not null;index"`
	CreatedByUserID uint `gorm:"not null;default:0;index"`

	RecurWeek    *int `gorm:"type:int"`
	RecurWeekday *int `gorm:"type:int"`

	Mentor  *User           `gorm:"foreignKey:MentorID"`
	Mentee  *User           `gorm:"foreignKey:MenteeID"`
	Mentees []PairingMentee `gorm:"foreignKey:PairingID"`

	Meetings []Meeting       `gorm:"foreignKey:PairingID"`
	Sharers  []PairingSharer `gorm:"foreignKey:PairingID"`
}
