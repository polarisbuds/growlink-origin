package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MeetingComment struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	Content string `gorm:"type:text;not null"`

	MeetingID uint `gorm:"not null;index"`
	UserID    uint `gorm:"not null;index"`

	Meeting *Meeting `gorm:"foreignKey:MeetingID"`
	User    *User    `gorm:"foreignKey:UserID"`
}
