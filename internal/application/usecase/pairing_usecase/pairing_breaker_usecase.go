package pairing_usecase

import (
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

func (uc *PairingBreakerUsecase) BreakPair(menteeID uint) bool {
	pairing, err := uc.pairingRepo.GetByMenteeID(menteeID)
	if err != nil {
		return false
	}
	if pairing == nil {
		return false
	}
	if err := uc.pairingRepo.BreakPairing(pairing.ID); err != nil {
		return false
	}
	return true
}
