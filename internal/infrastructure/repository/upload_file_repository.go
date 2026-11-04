package repository

import (
	"bytes"
	"context"
	"fmt"
	"growlink/internal/shared/utils"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
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
		return fmt.Errorf("[CREATE]アイコン保存先のフォルダの取得に失敗しました: %w", err)
	}

	filePath := filepath.Join(uploadsDirectory, name)
	if _, err := os.Stat(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("[CREATE]アイコンファイルの存在確認に失敗しました: %w", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ファイルの生成に失敗しました: %w", err)
	}
	defer f.Close()

	resizedImage, err := resizeImage(content)
	if err != nil {
		return fmt.Errorf("画像ファイルのリサイズに失敗しました: %w", err)
	}

	if _, err := io.Copy(f, resizedImage); err != nil {
		return fmt.Errorf("データの書き込みで失敗しました: %w", err)
	}
	return nil
}

func (repo *StorageFileRepository) DeleteFile(name string) error {
	uploadsDirectory, err := utils.GetAbsoluteTargetDir(repo.root)
	if err != nil {
		return fmt.Errorf("[DELETE]アイコン保存先のフォルダの取得に失敗しました: %w", err)
	}

	filePath := filepath.Join(uploadsDirectory, name)

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("[DELETE]アイコンファイルの存在確認に失敗しました: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("削除対象のPATHがディレクトリになっています")
	}

	if err = os.Remove(filePath); err != nil {
		return fmt.Errorf("アイコンファイルの削除に失敗しました: %w", err)
	}
	return nil
}

// 15MB
const (
	MaxIconUploadBytes = 15 << 20
	ResizedImageSize   = 128
)

func resizeImage(reader io.Reader) (io.Reader, error) {
	limited := io.LimitReader(reader, MaxIconUploadBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("アップロードファイルの読込に失敗しました: %w", err)
	}
	if len(data) > MaxIconUploadBytes {
		return nil, fmt.Errorf("画像が大きすぎます")
	}
	orgImg, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("イメージファイルのデコードに失敗しました: %w", err)
	}

	var (
		buf        = bytes.NewBuffer(nil)
		resizedImg = imaging.Fill(orgImg, ResizedImageSize, ResizedImageSize, imaging.Center, imaging.Lanczos)
	)

	if err := png.Encode(buf, resizedImg); err != nil {
		return nil, fmt.Errorf("リサイズ後イメージのエンコードに失敗しました: %w", err)
	}

	return bytes.NewReader(buf.Bytes()), nil
}
