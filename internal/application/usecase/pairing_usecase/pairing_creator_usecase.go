package pairing_usecase

import (
	"context"
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"
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

func (uc *PairingCreatorUsecase) PairingCreator(menteeID, mentorID uint) (err error) {
	const op = "pairing_usecase.PairingCreator"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	_, err = uc.pairingRepo.CreatePairing(menteeID, mentorID, 0)
	if err != nil {
		return fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
	}
	return nil
}
