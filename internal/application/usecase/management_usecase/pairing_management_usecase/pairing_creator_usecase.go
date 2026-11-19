package pairing_mng_usecase

import (
	"context"
	"fmt"

	inputmodel "growlink/internal/application/model/input_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"

	"github.com/google/uuid"
)

type PairingCreator struct {
	pairingRepo repository_interface.PairingRepository
	userRepo    repository_interface.UserRepository
}

func NewPairingCreator(pairingRepo repository_interface.PairingRepository, userRepo repository_interface.UserRepository) *PairingCreator {
	return &PairingCreator{
		pairingRepo: pairingRepo,
		userRepo:    userRepo,
	}
}

func (uc PairingCreator) CreatePairing(menteeID, mentorID uint) (err error) {
	const op = "pairing_mng_usecase.CreatePairing"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	return uc.CreatePairingFromInput(inputmodel.PairingInput{
		MentorID: mentorID,
		Targets:  []inputmodel.PairingTargetInput{{Kind: entity.PairingTargetKindUser, TargetID: menteeID}},
	}, 0)
}

func (uc PairingCreator) CreatePairingFromInput(input inputmodel.PairingInput, createdByUserID uint) (err error) {
	const op = "pairing_mng_usecase.CreatePairingFromInput"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	targets := input.ResolvedTargets()
	if len(targets) == 0 {
		return fmt.Errorf("対象が指定されていません")
	}

	mentorUUID, _ := uuid.Parse(input.MentorUuid)
	mentor, err := uc.userRepo.GetByUUID(mentorUUID)
	if err != nil {
		return fmt.Errorf("メンターの情報取得に失敗しました: %w", err)
	}

	menteeUUIDs := make([]string, len(targets))
	for idx, target := range targets {
		menteeUUIDs[idx] = target.TargetUuid
	}

	mentees, err := uc.userRepo.GetByUUIDs(menteeUUIDs)
	if err != nil {
		return fmt.Errorf("メンティーの情報取得に失敗しました: %w", err)
	}

	for _, mentee := range mentees {
		if _, err = uc.pairingRepo.CreatePairing(mentee.ID, mentor.ID, createdByUserID); err != nil {
			return fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
		}
	}
	return nil
}
