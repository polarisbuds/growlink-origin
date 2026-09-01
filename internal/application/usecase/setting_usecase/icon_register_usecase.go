package setting_usecase

import (
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
	"io"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type IconRegisterUsecase struct {
	uploadFileRepository repository_interface.UploadFileRepository
	userRepository       repository_interface.UserRepository
}

func NewIconRegisterUsecase(
	fileRepo repository_interface.UploadFileRepository,
	userRepo repository_interface.UserRepository) *IconRegisterUsecase {
	return &IconRegisterUsecase{
		uploadFileRepository: fileRepo,
		userRepository:       userRepo,
	}
}

func (uc *IconRegisterUsecase) RegistIcon(userID int64, name string, content io.Reader) error {
	user, err := uc.userRepository.GetByID(userID)
	if err != nil {
		return fmt.Errorf("ユーザアイコンのアップロードで存在しないユーザIDが指定されました: %w", err)
	}

	iconName := strings.Builder{}
	iconName.WriteString(uuid.NewString())
	ext := filepath.Ext(name)
	if !strings.HasPrefix(ext, ".") {
		iconName.WriteString(".")
	}
	iconName.WriteString(ext)
	filename := iconName.String()

	if err := uc.uploadFileRepository.NewFile(filename, content); err != nil {
		return fmt.Errorf("ユーザアイコンの新規作成処理に失敗しました: %w", err)
	}

	user.ImageURI = &filename
	if err := uc.userRepository.UpdateByID(user); err != nil {
		return fmt.Errorf("ユーザデータの更新（ユーザアイコン）に失敗しました: %w", err)
	}

	if filename := user.ImageURI; filename != nil {
		name := *filename
		slog.Info("古いユーザアイコンの削除を実施します。", slog.String("filename", name))

		isExists, err := uc.uploadFileRepository.Exists(name)
		if err != nil {
			return fmt.Errorf("アップロード済みのイメージファイルの存在チェックに失敗しました: %w", err)
		}
		if isExists {
			if err := uc.uploadFileRepository.DeleteFile(name); err != nil {
				return fmt.Errorf("アップロード済みのイメージファイルの削除に失敗しました: %w", err)
			}
		}
	}

	return nil
}
