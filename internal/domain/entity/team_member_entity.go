package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamMember struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	IsReader bool `gorm:"not null;default:false"`

	TeamID   int64 `gorm:"not null;index"`
	MemberID int64 `gorm:"not null;index"`

	Team   *Team `gorm:"foreignKey:TeamID"`
	Member *User `gorm:"foreignKey:MemberID"`
}
