package pairing_mng_usecase

import (
	"fmt"

	inputmodel "growlink/internal/application/model/input_model"
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
	return uc.CreatePairingFromInput(inputmodel.PairingInput{MenteeID: menteeID, MentorID: mentorID})
}

func (uc PairingCreator) CreatePairingFromInput(input inputmodel.PairingInput) error {
	pair, err := uc.pairingRepo.CreatePairing(input.MenteeID, input.MentorID)
	if err != nil {
		return fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
	}
	if err := uc.pairingRepo.ApplyPairingOptions(pair.ID, input.MentorID, input.MenteeID, input.RecurWeek, input.RecurWeekday, input.SharerUserIDs); err != nil {
		return err
	}
	return nil
}
