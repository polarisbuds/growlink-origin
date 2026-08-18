package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnswerComment struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	Content *string `gorm:"type:text"`

	AnswerID int64 `gorm:"not null;index"`
	UserID   int64 `gorm:"not null;index"`

	Answer *Answer `gorm:"foreignKey:AnswerID"`
	User   *User   `gorm:"foreignKey:UserID"`
}
