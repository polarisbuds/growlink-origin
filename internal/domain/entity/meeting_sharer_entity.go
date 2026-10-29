package entity

import "gorm.io/gorm"

type MeetingSharer struct {
	gorm.Model

	MeetingID uint `gorm:"not null;index"`
	UserID    uint `gorm:"not null;index"`

	User *User `gorm:"foreignKey:UserID"`
}
