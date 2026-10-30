package entity

import "gorm.io/gorm"

const (
	PairingTargetKindUser       = "user"
	PairingTargetKindTeam       = "team"
	PairingTargetKindDepartment = "department"
)

type PairingTarget struct {
	gorm.Model

	PairingID uint   `gorm:"not null;index"`
	Kind      string `gorm:"type:text;not null;index"`
	TargetID  uint   `gorm:"not null;index"`
}
