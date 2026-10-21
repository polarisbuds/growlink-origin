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

	RemindAutoSend bool   `gorm:"not null;default:false"`
	RemindWeekday  int    `gorm:"not null;default:1"`
	RemindTime     string `gorm:"type:text;not null;default:'09:00'"`

	DefaultTemplate *Template `gorm:"foreignKey:DefaultTemplateID"`
}
