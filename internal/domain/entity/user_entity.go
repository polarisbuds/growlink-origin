package entity

import (
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	UUID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Name        string    `gorm:"type:text;not null"`
	Email       string    `gorm:"type:text;not null"`
	Password    *string   `gorm:"type:text"`
	ImageURI    *string   `gorm:"type:text"`
	ImageBase64 *string   `gorm:"-"`

	Authority     valueobject.Authority  `gorm:"type:smallint;not null;index"`
	SystemRole    valueobject.SystemRole `gorm:"type:smallint;not null;default:1;index"`
	Suspended     bool                   `gorm:"not null;default:false;index"`
	LineWorksID   string                 `gorm:"type:text;not null;default:''"`
	RemindEnabled bool                   `gorm:"not null;default:false"`

	TeamMembers []TeamMember `gorm:"foreignKey:MemberID;references:ID"`

	DepartmentMembers []DepartmentMember `gorm:"foreignKey:UserID;references:ID"`
}
