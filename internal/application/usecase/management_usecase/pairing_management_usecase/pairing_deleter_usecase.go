package pairing_mng_usecase

import (
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
)

type PairingDeleter struct {
	pairingRepo repository_interface.PairingRepository
}

func NewPairingDeleter(pairingRepo repository_interface.PairingRepository) *PairingDeleter {
	return &PairingDeleter{
		pairingRepo: pairingRepo,
	}
}

func (uc *PairingDeleter) BreakPairing(pairingID uint) error {
	if err := uc.pairingRepo.BreakPairing(pairingID); err != nil {
		return fmt.Errorf("ペアリングの削除に失敗しました: %w", err)
	}
	return nil
}

func (uc *PairingDeleter) BreakPairings(pairingIDs ...uint) error {
	if err := uc.pairingRepo.BreakPairings(pairingIDs...); err != nil {
		return fmt.Errorf("複数ペアリングの削除に失敗しました: %w", err)
	}
	return nil
}
