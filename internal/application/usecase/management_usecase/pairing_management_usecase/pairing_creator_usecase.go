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
	return uc.CreatePairingFromInput(inputmodel.PairingInput{MenteeID: menteeID, MentorID: mentorID}, 0)
}

func (uc PairingCreator) CreatePairingFromInput(input inputmodel.PairingInput, createdByUserID uint) error {
	menteeIDs := input.ResolvedMenteeIDs()
	if len(menteeIDs) == 0 {
		return fmt.Errorf("メンティーが指定されていません")
	}
	pair, err := uc.pairingRepo.CreatePairing(menteeIDs[0], input.MentorID, createdByUserID)
	if err != nil {
		return fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
	}
	if err := uc.pairingRepo.ApplyPairingOptions(pair.ID, input.MentorID, menteeIDs, input.RecurWeek, input.RecurWeekday, input.SharerUserIDs); err != nil {
		return err
	}
	return nil
}
