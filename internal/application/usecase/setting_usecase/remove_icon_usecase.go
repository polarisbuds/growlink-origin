package setting_usecase

import (
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
	"log/slog"
)

type RemoveIconUsecase struct {
	uploadFileRepository repository_interface.UploadFileRepository
	userRepository       repository_interface.UserRepository
}

func NewRemoveIconUsecase(
	fileRepo repository_interface.UploadFileRepository,
	userRepo repository_interface.UserRepository) *RemoveIconUsecase {
	return &RemoveIconUsecase{
		uploadFileRepository: fileRepo,
		userRepository:       userRepo,
	}
}

func (uc *RemoveIconUsecase) RemoveIcon(userID int64) error {
	user, err := uc.userRepository.GetByID(userID)
	if err != nil {
		return fmt.Errorf("ユーザアイコンの削除処理で存在しないユーザIDが指定されました: %w", err)
	}

	if filename := user.ImageURI; filename != nil {
		name := *filename
		slog.Info("ユーザアイコンの削除を実施します。", slog.String("filename", name))

		isExists, err := uc.uploadFileRepository.Exists(name)
		if err != nil {
			return fmt.Errorf("削除対象のイメージファイルの存在チェックに失敗しました: %w", err)
		}
		if isExists {
			if err := uc.uploadFileRepository.DeleteFile(name); err != nil {
				return fmt.Errorf("削除対象のイメージファイルの削除に失敗しました: %w", err)
			}
		}
	}

	user.ImageURI = nil
	if err := uc.userRepository.UpdateByID(user); err != nil {
		return fmt.Errorf("ユーザデータの更新（ユーザアイコン）に失敗しました: %w", err)
	}

	return nil
}
