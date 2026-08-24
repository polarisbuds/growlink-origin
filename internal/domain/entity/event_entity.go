package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Name string    `gorm:"type:text;not null;uniqueIndex"`

	Histories []History `gorm:"foreignKey:EventID"`
}
