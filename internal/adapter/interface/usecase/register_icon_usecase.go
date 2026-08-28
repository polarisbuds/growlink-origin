package usecase_interface

import "io"

type SettingUsecase interface {
	IconUploader
}

type IconUploader interface {
	UploadIcon(userID int64, name string, content io.Reader) error
}
