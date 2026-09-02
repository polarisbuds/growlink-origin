package meeting_usecase

import (
	"fmt"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
)

type SubmitterUsecase struct {
	meetingRepo repository_interface.MeetingRepository
	answerRepo  repository_interface.AnswerRepository
}

func NewSubmitterUsecase(meetingRepo repository_interface.MeetingRepository, answerRepo repository_interface.AnswerRepository) *SubmitterUsecase {
	return &SubmitterUsecase{
		meetingRepo: meetingRepo,
		answerRepo:  answerRepo,
	}
}

func (uc *SubmitterUsecase) Submit(meetingID uint, answers []entity.Answer) error {
	if _, err := uc.answerRepo.Submits(answers); err != nil {
		return fmt.Errorf("回答の投稿に失敗しました: %w", err)
	}

	if err := uc.meetingRepo.UpdateStatus(meetingID, valueobject.MeetingStatusDone); err != nil {
		return fmt.Errorf("投稿完了へのステータス更新に失敗しました: %w", err)
	}
	return nil
}
