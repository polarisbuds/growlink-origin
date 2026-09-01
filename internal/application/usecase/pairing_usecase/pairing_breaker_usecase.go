package pairing_usecase

import (
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
)

type PairingBreakerUsecase struct {
	pairingRepo repository_interface.PairingRepository
	userRepo    repository_interface.UserRepository
}

func NewPairingBreakerUsecase(pairingRepo repository_interface.PairingRepository, userRepo repository_interface.UserRepository) *PairingBreakerUsecase {
	return &PairingBreakerUsecase{
		pairingRepo: pairingRepo,
		userRepo:    userRepo,
	}
}

func (uc *PairingBreakerUsecase) PairingCreator(menteeID uint) error {
	pairing, err := uc.pairingRepo.GetByMenteeID(menteeID)
	if err != nil {
		return fmt.Errorf("削除対象のペアリング情報の取得に失敗しました: %w", err)
	}

	if err := uc.pairingRepo.BreakPairing(pairing.ID); err != nil {
		return fmt.Errorf("ペアリングの削除に失敗しました: %w", err)
	}

	return nil
}
