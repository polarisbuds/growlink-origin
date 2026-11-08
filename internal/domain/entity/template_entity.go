package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Template struct {
	gorm.Model

	UUID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Title       string    `gorm:"type:text;not null;uniqueIndex"`
	Description string    `gorm:"type:text;not null;default:''"`

	DisplayNo int  `gorm:"not null;default:0;index"`
	IsSample  bool `gorm:"not null;default:false"`
	IsLock   bool `gorm:"not null;default:false"`

	UserID uint `gorm:"not null;index;column:creator_id"`

	User *User `gorm:"foreignKey:UserID"`

	Topics   []Topic   `gorm:"foreignKey:TemplateID"`
	Meetings []Meeting `gorm:"foreignKey:TemplateID"`
}
