package meeting_usecase

import (
	"fmt"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
)

type DraftSaverUsecase struct {
	meetingRepo repository_interface.MeetingRepository
	answerRepo  repository_interface.AnswerRepository
}

func NewDraftSaverUsecase(meetingRepo repository_interface.MeetingRepository, answerRepo repository_interface.AnswerRepository) *DraftSaverUsecase {
	return &DraftSaverUsecase{
		meetingRepo: meetingRepo,
		answerRepo:  answerRepo,
	}
}

func (uc *DraftSaverUsecase) SaveDraft(meetingID uint, answers []entity.Answer) error {
	if _, err := uc.answerRepo.Submits(answers); err != nil {
		return fmt.Errorf("回答の下書き保存に失敗しました: %w", err)
	}

	if err := uc.meetingRepo.UpdateStatus(meetingID, valueobject.MeetingStatusDraft); err != nil {
		return fmt.Errorf("下書き保存へのステータス更新に失敗しました: %w", err)
	}
	return nil
}
