package pairing_usecase

import (
	"fmt"
	"growlink/internal/domain/entity"
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

func (uc *PairingListCreatorUsecase) CreatePairList(mentorID uint) (menteeList []entity.User, err error) {
	pairings, err := uc.pairingRepo.GetByMentorID(mentorID)
	if err != nil {
		return nil, fmt.Errorf("ペアリストの作成に失敗しました: %w", err)
	}

	menties := []entity.User{}
	for _, pairing := range pairings {
		menties = append(menties, *pairing.Mentee)
	}

	return menties, nil
}
