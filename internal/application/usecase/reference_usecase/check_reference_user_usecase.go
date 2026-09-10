package referenceusecase

import (
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
)

type ReferenceCheckUsecase struct {
	UserReferenceUserRepository       repository_interface.UserReferenceUserRepository
	UserReferenceTeamRepository       repository_interface.UserReferenceTeamRepository
	UserReferenceDepartmentRepository repository_interface.UserReferenceDepartmentRepository
}

func NewReferenceCheck(UserReferenceUserRepository repository_interface.UserReferenceUserRepository,
	UserReferenceTeamRepository repository_interface.UserReferenceTeamRepository,
	UserReferenceDepartmentRepository repository_interface.UserReferenceDepartmentRepository) *ReferenceCheckUsecase {
	return &ReferenceCheckUsecase{
		UserReferenceUserRepository:       UserReferenceUserRepository,
		UserReferenceTeamRepository:       UserReferenceTeamRepository,
		UserReferenceDepartmentRepository: UserReferenceDepartmentRepository,
	}
}

func (uc *ReferenceCheckUsecase) CheckReferecneByUser(userID, targetUserID uint) (isOK bool, err error) {
	isReferable, err := uc.UserReferenceUserRepository.IsReferable(userID, targetUserID)
	if err != nil {
		return false, fmt.Errorf("ユーザID指定の参照可能チェックの処理に失敗しました: %w", err)
	}
	return isReferable, nil
}

func (uc *ReferenceCheckUsecase) CheckReferecneByTeam(userID, teamID uint) (isOK bool, err error) {
	isReferable, err := uc.UserReferenceTeamRepository.IsReferable(userID, teamID)
	if err != nil {
		return false, fmt.Errorf("チームID指定の参照可能チェックの処理に失敗しました: %w", err)
	}
	return isReferable, nil
}

func (uc *ReferenceCheckUsecase) CheckReferecneByDepartment(userID, departmentID uint) (isOK bool, err error) {
	isReferable, err := uc.UserReferenceDepartmentRepository.IsReferable(userID, departmentID)
	if err != nil {
		return false, fmt.Errorf("部署ID指定の参照可能チェックの処理に失敗しました: %w", err)
	}
	return isReferable, nil
}
