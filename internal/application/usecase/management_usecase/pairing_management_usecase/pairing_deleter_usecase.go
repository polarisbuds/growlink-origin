package pairing_mng_usecase

import (
	"context"
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"
)

type PairingDeleter struct {
	pairingRepo repository_interface.PairingRepository
}

func NewPairingDeleter(pairingRepo repository_interface.PairingRepository) *PairingDeleter {
	return &PairingDeleter{
		pairingRepo: pairingRepo,
	}
}

func (uc *PairingDeleter) BreakPairing(pairingID uint) (err error) {
	const op = "pairing_mng_usecase.BreakPairing"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.pairingRepo.BreakPairing(pairingID); err != nil {
		return fmt.Errorf("ペアリングの削除に失敗しました: %w", err)
	}
	return nil
}

func (uc *PairingDeleter) BreakPairings(pairingIDs ...uint) (err error) {
	const op = "pairing_mng_usecase.BreakPairings"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.pairingRepo.BreakPairings(pairingIDs...); err != nil {
		return fmt.Errorf("複数ペアリングの削除に失敗しました: %w", err)
	}
	return nil
}
