package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Template struct {
	gorm.Model

	UUID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Title string    `gorm:"type:text;not null;uniqueIndex"`

	UserID uint `gorm:"not null;index"`

	User *User `gorm:"foreignKey:UserID"`

	Topics   []Topic   `gorm:"foreignKey:TemplateID"`
	Meetings []Meeting `gorm:"foreignKey:TemplateID"`
}
