package usecase_interface

//go:generate mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import "io"

type SettingUsecase interface {
	IconUploader
	IconRemover
}

type IconUploader interface {
	UploadIcon(userID int64, name string, content io.Reader) error
}

type IconRemover interface {
	RemoveIcon(userID int64) error
}
