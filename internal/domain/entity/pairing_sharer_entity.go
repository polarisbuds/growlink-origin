package entity

import "gorm.io/gorm"

type PairingSharer struct {
	gorm.Model

	PairingID uint `gorm:"not null;index"`
	UserID    uint `gorm:"not null;index"`

	Pairing *Pairing `gorm:"foreignKey:PairingID"`
	User    *User    `gorm:"foreignKey:UserID"`
}
