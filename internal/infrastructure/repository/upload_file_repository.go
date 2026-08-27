package repository

import (
	"context"
	"fmt"
	"growlink/internal/shared/utils"
	"io"
	"os"
	"path/filepath"
)

const FIRST_GENERATION = 1

type StorageFileRepository struct {
	ctx  context.Context
	root string
}

func NewStorageRepository(ctx context.Context, root string) *StorageFileRepository {
	return &StorageFileRepository{
		ctx:  ctx,
		root: root,
	}
}

func (repo *StorageFileRepository) Exists(name string) (isExists bool, err error) {
	uploadsDirectory, err := utils.GetAbsoluteTargetDir(repo.root)
	filePath := filepath.Join(uploadsDirectory, name)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("ファイルの存在確認でエラーが発生しました: %w", err)
	}
	return true, nil
}

func (repo *StorageFileRepository) NewFile(name string, content io.Reader) error {
	uploadsDirectory, err := utils.GetAbsoluteTargetDir(repo.root)
	if err != nil {
		return fmt.Errorf("アイコン保存先のフォルダの取得に失敗しました: %w", err)
	}

	filePath := filepath.Join(uploadsDirectory, name)
	if _, err := os.Stat(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("ファイルの存在確認でエラーが発生しました: %w", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ファイルの生成に失敗しました: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, content); err != nil {
		return fmt.Errorf("データの書き込みで失敗しました: %w", err)
	}
	return nil
}

func (repo *StorageFileRepository) DeleteFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("削除対象のPATHがディレクトリになっています")
	}
	return os.Remove(path)
}
