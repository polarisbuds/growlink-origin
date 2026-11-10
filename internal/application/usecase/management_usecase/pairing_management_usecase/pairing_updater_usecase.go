package pairing_mng_usecase

import (
	"fmt"

	inputmodel "growlink/internal/application/model/input_model"
	"growlink/internal/domain/entity"
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
	menteeIDs := make([]uint, 0, len(targets))
	for _, t := range targets {
		if t.Kind != entity.PairingTargetKindUser {
			return fmt.Errorf("ユーザ以外の対象はサポートされていません")
		}
		menteeIDs = append(menteeIDs, t.TargetID)
	}
	existing, err := uc.pairingRepo.Get(pairingID)
	if err != nil {
		return fmt.Errorf("ペアリング情報の取得に失敗しました: %w", err)
	}
	if err := uc.pairingRepo.UpdatePairing(pairingID, input.MentorID, menteeIDs[0]); err != nil {
		return fmt.Errorf("ペアリング情報の更新に失敗しました: %w", err)
	}
	for _, menteeID := range menteeIDs[1:] {
		if _, err := uc.pairingRepo.CreatePairing(menteeID, input.MentorID, existing.CreatedByUserID); err != nil {
			return fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
		}
	}
	return nil
}
