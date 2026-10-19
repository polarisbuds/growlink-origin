package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type PairingRepository interface {
	CreatePairing(menteeID, mentorID uint) (*entity.Pairing, error)
	BreakPairing(pairingID uint) error
	BreakPairings(pairingID ...uint) error
	BreakAllMenteePairing(mentorID uint) error
	GetAll() ([]entity.Pairing, error)
	Get(pairingID uint) (*entity.Pairing, error)
	GetByMenteeID(menteeID uint) (*entity.Pairing, error)
	GetByMentorID(mentorID uint) ([]entity.Pairing, error)
	UpdatePairing(pairingID, mentorID uint) error
	ApplyPairingOptions(pairingID uint, mentorID, menteeID uint, recurWeek, recurWeekday *int, sharerUserIDs []uint) error
}
