package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type History struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	UserID  int64 `gorm:"not null;index"`
	EventID int64 `gorm:"not null;index"`

	User  *User  `gorm:"foreignKey:UserID"`
	Event *Event `gorm:"foreignKey:EventID"`
}
