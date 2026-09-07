package usecase_interface

import (
	outputmodel "growlink/internal/application/model/output_model"
)

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type (
	PairingUsecase interface {
		PairingListCreator
		PairingBreaker
	}

	PairingListCreator interface {
		CreatePairList(mentorID uint) (menteeList []outputmodel.UserOutput, err error)
	}
	PairingBreaker interface {
		BreakPair(menteeID uint) bool
	}
)
