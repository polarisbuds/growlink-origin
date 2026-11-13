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
	_, err := uc.pairingRepo.CreatePairing(menteeID, mentorID, 0)
	if err != nil {
		return fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
	}
	return nil
}
