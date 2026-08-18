package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Topic struct {
	gorm.Model

	UUID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	DisplayNo int64     `gorm:"not null"`
	Question  string    `gorm:"type:text;not null"`
	IsRequire bool      `gorm:"not null"`

	TemplateID int64 `gorm:"not null;index"`

	Template *Template `gorm:"foreignKey:TemplateID"`

	Answers []Answer `gorm:"foreignKey:TopicID"`
}
