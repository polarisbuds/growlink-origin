package pairing_usecase

import (
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
)

type PairingListCreatorUsecase struct {
	pairingRepo repository_interface.PairingRepository
	userRepo    repository_interface.UserRepository
}

func NewPairingListCreatorUsecase(pairingRepo repository_interface.PairingRepository, userRepo repository_interface.UserRepository) *PairingListCreatorUsecase {
	return &PairingListCreatorUsecase{
		pairingRepo: pairingRepo,
		userRepo:    userRepo,
	}
}

func (uc *PairingListCreatorUsecase) CreatePairList(mentorID uint) (menteeList []outputmodel.UserOutput, err error) {
	pairings, err := uc.pairingRepo.GetByMentorID(mentorID)
	if err != nil {
		return nil, fmt.Errorf("ペアリストの作成に失敗しました: %w", err)
	}

	menties := []outputmodel.UserOutput{}
	for _, pairing := range pairings {
		if pairing.Mentee != nil {
			menties = append(menties, converter.UserEntityToOutput(*pairing.Mentee))
		}
	}
	return menties, nil
}
