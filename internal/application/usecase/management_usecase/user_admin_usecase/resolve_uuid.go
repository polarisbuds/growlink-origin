package user_admin_usecase

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (uc *UserAdmin) ResolveUserID(userUUID string) (uint, error) {
	parsed, err := uuid.Parse(userUUID)
	if err != nil {
		return 0, err
	}
	user, err := uc.userRepo.GetByUUID(parsed)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (uc *UserAdmin) UpdateUserByUUID(userUUID string, input UpdateUserInput) error {
	userID, err := uc.ResolveUserID(userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return fmt.Errorf("invalid user uuid")
	}
	return uc.UpdateUser(userID, input)
}