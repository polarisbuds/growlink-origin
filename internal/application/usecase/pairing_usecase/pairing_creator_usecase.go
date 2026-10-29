package pairing_usecase

import (
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
)

type PairingCreatorUsecase struct {
	pairingRepo repository_interface.PairingRepository
	userRepo    repository_interface.UserRepository
}

func NewPairingCreatorUsecase(pairingRepo repository_interface.PairingRepository, userRepo repository_interface.UserRepository) *PairingCreatorUsecase {
	return &PairingCreatorUsecase{
		pairingRepo: pairingRepo,
		userRepo:    userRepo,
	}
}

func (uc *PairingCreatorUsecase) PairingCreator(menteeID, mentorID uint) error {
	oldPair, err := uc.pairingRepo.GetByMenteeID(menteeID)
	if err != nil {
		return fmt.Errorf("ペアリング（旧）の取得に失敗しました: %w", err)
	}

	_, err = uc.pairingRepo.CreatePairing(menteeID, mentorID, 0)
	if err != nil {
		return fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
	}

	if oldPair != nil {
		if err := uc.pairingRepo.BreakPairing(oldPair.ID); err != nil {
			return fmt.Errorf("ペアリングの解消に失敗しました: %w", err)
		}
	}

	return nil
}
