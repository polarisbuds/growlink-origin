package usecase_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type PairingUsecase interface {
	PairingListCreator
	PairingCreator
	PairingBreaker
}

type PairingListCreator interface {
	CreatePairList(mentorID int64) (menteeList []entity.User, err error)
}

type PairingCreator interface {
	CreatePair(menteeID, mentorID int64) bool
}

type PairingBreaker interface {
	BreakPair(menteeID, mentorID int64) bool
}
