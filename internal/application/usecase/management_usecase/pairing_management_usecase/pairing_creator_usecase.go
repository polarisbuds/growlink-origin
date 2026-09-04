package pairing_mng_usecase

import (
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
)

type PairingCreator struct {
	pairingRepo repository_interface.PairingRepository
}

func NewPairingCreator(pairingRepo repository_interface.PairingRepository) *PairingCreator {
	return &PairingCreator{
		pairingRepo: pairingRepo,
	}
}

func (uc PairingCreator) CreatePairing(menteeID, mentorID uint) error {
	if _, err := uc.pairingRepo.CreatePairing(menteeID, mentorID); err != nil {
		return fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
	}
	return nil
}
