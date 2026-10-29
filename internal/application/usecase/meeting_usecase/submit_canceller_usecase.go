package meeting_usecase

import (
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
)

type SubmitCancellerUsecase struct {
	meetingRepo repository_interface.MeetingRepository
}

func NewSubmitCancellerUsecase(meetingRepo repository_interface.MeetingRepository) *SubmitCancellerUsecase {
	return &SubmitCancellerUsecase{
		meetingRepo: meetingRepo,
	}
}

func (uc *SubmitCancellerUsecase) CancelSubmit(meetingID uint) error {
	if err := uc.meetingRepo.UpdateStatus(meetingID, valueobject.MeetingStatusDraft); err != nil {
		return fmt.Errorf("下書き保存へのステータス更新に失敗しました: %w", err)
	}
	return nil
}
