package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SystemSettings struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	MeetingWeek       int    `gorm:"not null"`
	MeetingWeekday    int    `gorm:"not null"`
	DefaultTemplateID uint   `gorm:"not null;index"`
	LogOutputPath     string `gorm:"type:text"`
	LogFilePrefix     string `gorm:"type:text"`

	DefaultTemplate *Template `gorm:"foreignKey:DefaultTemplateID"`
}
