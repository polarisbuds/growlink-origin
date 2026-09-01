package usecase_interface

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import "io"

type SettingUsecase interface {
	IconRegister
	IconRemover
}

type IconRegister interface {
	RegistIcon(userID int64, name string, content io.Reader) error
}

type IconRemover interface {
	RemoveIcon(userID int64) error
}
