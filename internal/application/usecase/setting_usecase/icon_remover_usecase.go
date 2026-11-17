package setting_usecase

import (
	"context"
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"
	"log/slog"
)

type IconRemoverUsecase struct {
	uploadFileRepository repository_interface.UploadFileRepository
	userRepository       repository_interface.UserRepository
}

func NewIconRemoverUsecase(
	fileRepo repository_interface.UploadFileRepository,
	userRepo repository_interface.UserRepository) *IconRemoverUsecase {
	return &IconRemoverUsecase{
		uploadFileRepository: fileRepo,
		userRepository:       userRepo,
	}
}

func (uc *IconRemoverUsecase) RemoveIcon(userID uint) (err error) {
	const op = "setting_usecase.RemoveIcon"
	ctx := context.Background()
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	user, err := uc.userRepository.GetByID(userID)
	if err != nil {
		return fmt.Errorf("ユーザアイコンの削除処理で存在しないユーザIDが指定されました: %w", err)
	}

	if filename := user.ImageURI; filename != nil {
		name := *filename
		logutils.Info(ctx, logutils.LayerUsecase, op, "ユーザアイコンの削除を実施します。", slog.String("filename", name))

		isExists, err := uc.uploadFileRepository.Exists(name)
		if err != nil {
			return fmt.Errorf("削除対象のイメージファイルの存在チェックに失敗しました: %w", err)
		}
		if isExists {
			if err = uc.uploadFileRepository.DeleteFile(name); err != nil {
				return fmt.Errorf("削除対象のイメージファイルの削除に失敗しました: %w", err)
			}
		}
	}

	if err = uc.userRepository.UpdateColumns(userID, map[string]any{"image_uri": nil}); err != nil {
		return fmt.Errorf("ユーザアイコンの削除に失敗しました: %w", err)
	}

	return nil
}
