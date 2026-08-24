package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	UUID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`
	Name     string    `gorm:"type:text;not null"`
	Email    string    `gorm:"type:text;not null"`
	Password *string   `gorm:"type:text"`
	ImageURI *string   `gorm:"type:text"`

	RoleID       int64 `gorm:"not null;index"`
	DepartmentID int64 `gorm:"not null;index"`
	AuthorityID  int64 `gorm:"not null;index"`

	Role       *Role       `gorm:"foreignKey:RoleID"`
	Department *Department `gorm:"foreignKey:DepartmentID"`
	Authority  *Authority  `gorm:"foreignKey:AuthorityID"`

	TeamMembers                  []TeamMember       `gorm:"foreignKey:MemberID;references:ID"`
	UserRelationshipsSourceUsers []UserRelationship `gorm:"foreignKey:SourceUserID;references:ID"`
	UserRelationshipsTargetUsers []UserRelationship `gorm:"foreignKey:TargetUserID;references:ID"`
}
