package pairing_mng_usecase

import (
	"context"
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	logutils "growlink/internal/shared/logUtils"
)

func (uc *PairingUpdater) ListSharers(pairingID uint) (ret0 []outputmodel.UserOutput, err error) {
	const op = "pairing_mng_usecase.ListSharers"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	sharers, err := uc.pairingRepo.ListSharersByPairingID(pairingID)
	if err != nil {
		return nil, err
	}
	return converter.PairingSharersToUserOutputs(sharers), nil
}

func (uc *PairingUpdater) ReplaceSharers(pairingID uint, userIDs []uint) (err error) {
	const op = "pairing_mng_usecase.ReplaceSharers"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.pairingRepo.ReplaceSharers(pairingID, userIDs); err != nil {
		return fmt.Errorf("ペア共有者の更新に失敗しました: %w", err)
	}
	return nil
}