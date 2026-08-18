package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MeetingComment struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	Content *string `gorm:"type:text"`

	MeetingID int64 `gorm:"not null;index"`
	UserID    int64 `gorm:"not null;index"`

	Meeting *Meeting `gorm:"foreignKey:MeetingID"`
	User    *User    `gorm:"foreignKey:UserID"`
}
