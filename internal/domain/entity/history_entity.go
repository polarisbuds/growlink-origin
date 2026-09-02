package entity

import (
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type History struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	UserID    uint                  `gorm:"not null;index"`
	EventType valueobject.EventType `gorm:"type:text;not null;index"`

	User *User `gorm:"foreignKey:UserID"`
}
