package repository_interface

//go:generate mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import "io"

type UploadFileRepository interface {
	// ファイルの作成
	Exists(name string) (isExists bool, err error)
	// ファイルの作成
	NewFile(name string, content io.Reader) error
	// ファイルの削除
	DeleteFile(path string) error
}
