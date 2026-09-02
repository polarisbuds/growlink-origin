package usecase_interface

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import "io"

type (
	SettingUsecase interface {
		IconRegister
		IconRemover
	}

	IconRegister interface {
		RegistIcon(userID uint, name string, content io.Reader) error
	}
	IconRemover interface {
		RemoveIcon(userID uint) error
	}
)
