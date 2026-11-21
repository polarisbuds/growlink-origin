package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LineWorksSetting struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	Title   string `gorm:"type:text;not null;uniqueIndex:idx_line_works_setting_title,where:deleted_at IS NULL;uniqueIndex:idx_line_works_setting"`
	Content string `gorm:"type:text;not null;uniqueIndex:idx_line_works_setting"`

	UpdateUserID uint `gorm:"not null;default:0;index"`
}