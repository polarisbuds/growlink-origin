package setting_usecase

import (
	"context"

	user_admin_usecase "growlink/internal/application/usecase/management_usecase/user_admin_usecase"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"

	"github.com/google/uuid"
)

type ProfileOutput struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	LineWorksID     string  `json:"lineWorksId"`
	ImageURI      *string `json:"imageUri"`
	RemindEnabled bool    `json:"remindEnabled"`
}

type ProfileUsecase struct {
	userRepo repository_interface.UserRepository
}

func NewProfileUsecase(userRepo repository_interface.UserRepository) *ProfileUsecase {
	return &ProfileUsecase{userRepo: userRepo}
}

func (uc *ProfileUsecase) GetProfile(userID uint) (ret0 *ProfileOutput, err error) {
	const op = "setting_usecase.GetProfile"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() { done(err) }()
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return &ProfileOutput{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		LineWorksID:   user.LineWorksID,
		ImageURI:      user.ImageURI,
		RemindEnabled: user.RemindEnabled,
	}, nil
}

func (uc *ProfileUsecase) UpdateProfile(userID uint, lineWorksID string, remindEnabled *bool) error {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if err := user_admin_usecase.ValidateLineWorksLocalID(lineWorksID); err != nil {
		return err
	}
	user.LineWorksID = lineWorksID
	if remindEnabled != nil {
		user.RemindEnabled = *remindEnabled
	}
	return uc.userRepo.UpdateByID(user)
}

func (uc *ProfileUsecase) GetUserWithAvatarByID(userID uint) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	uc.userRepo.AttachImageBase64(user)
	return user, nil
}

func (uc *ProfileUsecase) GetUserWithAvatarByUUID(userUUID string) (*entity.User, error) {
	parsed, err := uuid.Parse(userUUID)
	if err != nil {
		return nil, err
	}
	user, err := uc.userRepo.GetByUUID(parsed)
	if err != nil {
		return nil, err
	}
	uc.userRepo.AttachImageBase64(user)
	return user, nil
}