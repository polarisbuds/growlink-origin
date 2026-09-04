package pairing_mng_usecase

import (
	"fmt"
	inputmodel "growlink/internal/application/model/input_model"
	repository_interface "growlink/internal/domain/interface/repository"
)

type PairingUpdater struct {
	pairingRepo repository_interface.PairingRepository
}

func NewPairingUpdater(pairingRepo repository_interface.PairingRepository) *PairingUpdater {
	return &PairingUpdater{
		pairingRepo: pairingRepo,
	}
}

func (uc *PairingUpdater) UpdatePairing(pairingID uint, input inputmodel.PairingInput) error {
	if err := uc.pairingRepo.UpdatePairing(pairingID, input.MentorID); err != nil {
		return fmt.Errorf("ペアリング情報の更新に失敗しました: %w", err)
	}
	return nil
}
