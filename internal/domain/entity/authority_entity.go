package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Authority struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Type string    `gorm:"type:varchar(16);not null;uniqueIndex"`

	Users []User `gorm:"foreignKey:AuthorityID"`
}
