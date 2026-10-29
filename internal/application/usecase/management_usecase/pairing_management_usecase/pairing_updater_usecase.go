package pairing_mng_usecase

import (
	"fmt"

	"growlink/internal/application/converter"
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
	targets := input.ResolvedTargets()
	if len(targets) == 0 {
		return fmt.Errorf("対象が指定されていません")
	}
	opts := converter.PairingTargetInputsToOptions(targets)
	if err := uc.pairingRepo.ApplyPairingOptions(pairingID, input.MentorID, opts); err != nil {
		return fmt.Errorf("ペアリング情報の更新に失敗しました: %w", err)
	}
	return nil
}
