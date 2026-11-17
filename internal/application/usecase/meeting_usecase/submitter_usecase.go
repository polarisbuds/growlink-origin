package meeting_usecase

import (
	"context"
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"
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

func (uc *SubmitterUsecase) Submit(meetingID uint, answers []outputmodel.AnswerOutput, nextStatus valueobject.MeetingStatus) (err error) {
	const op = "meeting_usecase.Submit"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if _, err = uc.answerRepo.Submits(converter.AnswerOutputsToEntities(answers)); err != nil {
		return fmt.Errorf("回答の投稿に失敗しました: %w", err)
	}

	if err = uc.meetingRepo.UpdateStatus(meetingID, nextStatus); err != nil {
		return fmt.Errorf("投稿完了へのステータス更新に失敗しました: %w", err)
	}
	return nil
}
