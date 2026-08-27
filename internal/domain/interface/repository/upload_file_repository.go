package repository_interface

import "io"

type UploadFileRepository interface {
	// ファイルの作成
	Exists(name string) (isExists bool, err error)
	// ファイルの作成
	NewFile(name string, content io.Reader) error
	// ファイルの削除
	DeleteFile(path string) error
}
