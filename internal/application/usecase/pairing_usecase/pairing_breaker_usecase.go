package pairing_usecase

import (
	"context"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/shared/logUtils"
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

func (uc *PairingBreakerUsecase) BreakPairByMenteeID(menteeID uint) bool {
	const op = "pairing_usecase.BreakPairByMenteeID"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(nil)
	}()
	pairings, err := uc.pairingRepo.GetByMenteeID(menteeID)
	if err != nil {
		return false
	}
	if len(pairings) == 0 {
		return false
	}
	for _, pairing := range pairings {
		if err := uc.pairingRepo.BreakPairing(pairing.ID); err != nil {
			return false
		}
	}
	return true
}
