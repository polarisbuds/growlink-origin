package pairing_mng_usecase

import (
	"fmt"

	inputmodel "growlink/internal/application/model/input_model"
	"growlink/internal/domain/entity"
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
	return uc.CreatePairingFromInput(inputmodel.PairingInput{
		MentorID: mentorID,
		Targets:  []inputmodel.PairingTargetInput{{Kind: entity.PairingTargetKindUser, TargetID: menteeID}},
	}, 0)
}

func (uc PairingCreator) CreatePairingFromInput(input inputmodel.PairingInput, createdByUserID uint) error {
	targets := input.ResolvedTargets()
	if len(targets) == 0 {
		return fmt.Errorf("対象が指定されていません")
	}
	for _, t := range targets {
		if t.Kind != entity.PairingTargetKindUser {
			return fmt.Errorf("ユーザ以外の対象はサポートされていません")
		}
		if _, err := uc.pairingRepo.CreatePairing(t.TargetID, input.MentorID, createdByUserID); err != nil {
			return fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
		}
	}
	return nil
}
