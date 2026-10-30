package meeting_usecase

import (
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
)

type AnswerReaderUsecase struct {
	answerRepo repository_interface.AnswerRepository
}

func NewAnswerReaderUsecase(answerRepo repository_interface.AnswerRepository) *AnswerReaderUsecase {
	return &AnswerReaderUsecase{
		answerRepo: answerRepo,
	}
}

func (uc *AnswerReaderUsecase) GetAnswer(answerID uint) (*outputmodel.AnswerOutput, error) {
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

func (uc *AnswerReaderUsecase) ListAnswersByMeeting(meetingID uint) ([]outputmodel.AnswerOutput, error) {
	answers, err := uc.answerRepo.GetByMeetingID(meetingID)
	if err != nil {
		return nil, fmt.Errorf("回答情報の取得に失敗しました: %w", err)
	}
	return converter.AnswerEntitiesToOutputs(answers), nil
}
