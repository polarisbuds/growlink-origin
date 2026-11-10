package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Template struct {
	gorm.Model

	UUID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Title       string    `gorm:"type:text;not null"`
	Description string    `gorm:"type:text;not null;default:''"`

	IsSample    bool `gorm:"not null;default:false"`
	IsLock      bool `gorm:"not null;default:false"`
	IsDefault   bool `gorm:"not null;default:false;column:is_default"`

	UserID uint `gorm:"not null;index;column:creator_id"`

	User *User `gorm:"foreignKey:UserID"`

	Topics   []Topic   `gorm:"foreignKey:TemplateID"`
	Meetings []Meeting `gorm:"foreignKey:TemplateID"`
}
