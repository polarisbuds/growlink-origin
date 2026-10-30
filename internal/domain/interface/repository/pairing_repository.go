package repository_interface

import (
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
)

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type PairingRepository interface {
	CreatePairing(menteeID, mentorID, createdByUserID uint) (*entity.Pairing, error)
	CreateTemporaryPairing(menteeID, mentorID, createdByUserID uint) (*entity.Pairing, error)
	BreakPairing(pairingID uint) error
	BreakPairings(pairingID ...uint) error
	BreakAllMenteePairing(mentorID uint) error
	GetAll() ([]entity.Pairing, error)
	Get(pairingID uint) (*entity.Pairing, error)
	GetByUUID(pairingUUID uuid.UUID) (*entity.Pairing, error)
	GetByMenteeID(menteeID uint) (*entity.Pairing, error)
	GetByMentorID(mentorID uint) ([]entity.Pairing, error)
	GetPairingIDsBySharerUserID(userID uint) ([]uint, error)
	UpdatePairing(pairingID, mentorID uint) error
	ApplyPairingOptions(pairingID uint, mentorID uint, targets []PairingTargetOption) error
}

type PairingTargetOption struct {
	Kind     string
	TargetID uint
}
