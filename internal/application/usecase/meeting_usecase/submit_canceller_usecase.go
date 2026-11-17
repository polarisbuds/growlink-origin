package meeting_usecase

import (
	"context"
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"
)

type SubmitCancellerUsecase struct {
	meetingRepo repository_interface.MeetingRepository
}

func NewSubmitCancellerUsecase(meetingRepo repository_interface.MeetingRepository) *SubmitCancellerUsecase {
	return &SubmitCancellerUsecase{
		meetingRepo: meetingRepo,
	}
}

func (uc *SubmitCancellerUsecase) CancelSubmit(meetingID uint) (err error) {
	const op = "meeting_usecase.CancelSubmit"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.meetingRepo.UpdateStatus(meetingID, valueobject.MeetingStatusDraft); err != nil {
		return fmt.Errorf("下書き保存へのステータス更新に失敗しました: %w", err)
	}
	return nil
}
