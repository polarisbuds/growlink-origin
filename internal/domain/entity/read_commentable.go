package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReadCommentable struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	MeetingID int64 `gorm:"not null;index"`
	UserID    int64 `gorm:"not null;index"`

	Meeting *Meeting `gorm:"foreignKey:MeetingID"`
	User    *User    `gorm:"foreignKey:UserID"`
}
