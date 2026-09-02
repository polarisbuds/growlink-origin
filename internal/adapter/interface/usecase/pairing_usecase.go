package usecase_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type (
	PairingUsecase interface {
		PairingListCreator
		PairingBreaker
	}

	PairingListCreator interface {
		CreatePairList(mentorID uint) (menteeList []entity.User, err error)
	}
	PairingBreaker interface {
		BreakPair(menteeID uint) bool
	}
)
