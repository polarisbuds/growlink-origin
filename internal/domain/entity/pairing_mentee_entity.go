package entity

import "gorm.io/gorm"

type PairingMentee struct {
	gorm.Model
	PairingID uint `gorm:"not null;index"`
	UserID    uint `gorm:"not null;index"`

	User *User `gorm:"foreignKey:UserID"`
}
