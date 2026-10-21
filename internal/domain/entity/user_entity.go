package entity

import (
	"growlink/internal/domain/valueobject"

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

	RoleID       uint                  `gorm:"not null;index"`
	DepartmentID uint                  `gorm:"not null;index"`
	Authority       valueobject.Authority       `gorm:"type:smallint;not null;index"`
	SystemRole      valueobject.SystemRole      `gorm:"type:smallint;not null;default:1;index"`
	OperationTarget valueobject.OperationTarget `gorm:"type:smallint;not null;default:1;index"`
	Suspended       bool                        `gorm:"not null;default:false;index"`
	LineWorksID     string                      `gorm:"type:text;not null;default:''"`
	RemindEnabled   bool                        `gorm:"not null;default:false"`

	Role       *Role       `gorm:"foreignKey:RoleID"`
	Department *Department `gorm:"foreignKey:DepartmentID"`

	TeamMembers                  []TeamMember       `gorm:"foreignKey:MemberID;references:ID"`
	UserRelationshipsSourceUsers []UserRelationship `gorm:"foreignKey:SourceUserID;references:ID"`
	UserRelationshipsTargetUsers []UserRelationship `gorm:"foreignKey:TargetUserID;references:ID"`

	DepartmentRoles              []DepartmentRole        `gorm:"foreignKey:UserID;references:ID"`
	UserReferenceUsers           []UserReferenceUser     `gorm:"foreignKey:UserID;references:ID"`
	ReferableUserReferences      []UserReferenceUser     `gorm:"foreignKey:ReferableUserID;references:ID"`
	UserReferenceTeams           []UserReferenceTeam     `gorm:"foreignKey:UserID;references:ID"`
	UserReferenceDepartments     []UserReferenceDepartment `gorm:"foreignKey:UserID;references:ID"`
}
