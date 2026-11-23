package entity

import (
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LineWorksMessageSetting struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	Kind valueobject.LineWorksMessageKind `gorm:"type:smallint;not null;uniqueIndex:idx_line_works_message_setting_kind,where:deleted_at IS NULL"`

	Title   string `gorm:"type:text;not null;uniqueIndex:idx_line_works_message_setting_title,where:deleted_at IS NULL;uniqueIndex:idx_line_works_message_setting"`
	Content string `gorm:"type:text;not null;uniqueIndex:idx_line_works_message_setting"`

	UpdateUserID uint `gorm:"not null;default:0;index"`
}

func (LineWorksMessageSetting) TableName() string {
	return "line_works_message_settings"
}