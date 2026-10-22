package pairing_mng_usecase

import (
	"fmt"

	"growlink/internal/application/converter"
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
	opts := converter.PairingTargetInputsToOptions(targets)
	primaryMentee := opts[0].TargetID
	for _, o := range opts {
		if o.Kind == entity.PairingTargetKindUser {
			primaryMentee = o.TargetID
			break
		}
	}
	pair, err := uc.pairingRepo.CreatePairing(primaryMentee, input.MentorID, createdByUserID)
	if err != nil {
		return fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
	}
	if err := uc.pairingRepo.ApplyPairingOptions(pair.ID, input.MentorID, opts); err != nil {
		return err
	}
	return nil
}
