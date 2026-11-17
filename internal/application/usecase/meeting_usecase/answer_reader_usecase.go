package meeting_usecase

import (
	"context"
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"
)

type AnswerReaderUsecase struct {
	answerRepo repository_interface.AnswerRepository
}

func NewAnswerReaderUsecase(answerRepo repository_interface.AnswerRepository) *AnswerReaderUsecase {
	return &AnswerReaderUsecase{
		answerRepo: answerRepo,
	}
}

func (uc *AnswerReaderUsecase) GetAnswer(answerID uint) (ret0 *outputmodel.AnswerOutput, err error) {
	const op = "meeting_usecase.GetAnswer"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	answer, err := uc.answerRepo.GetByID(answerID)
	if err != nil {
		return nil, fmt.Errorf("回答情報の取得に失敗しました: %w", err)
	}
	if answer == nil {
		return nil, fmt.Errorf("回答情報の取得に失敗しました: 指定された回答が存在しません")
	}
	out := converter.AnswerEntityToOutput(*answer)
	return &out, nil
}

func (uc *AnswerReaderUsecase) ListAnswersByMeeting(meetingID uint) (ret0 []outputmodel.AnswerOutput, err error) {
	const op = "meeting_usecase.ListAnswersByMeeting"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	answers, err := uc.answerRepo.GetByMeetingID(meetingID)
	if err != nil {
		return nil, fmt.Errorf("回答情報の取得に失敗しました: %w", err)
	}
	return converter.AnswerEntitiesToOutputs(answers), nil
}
