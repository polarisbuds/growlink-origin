package meeting_usecase

import (
	"fmt"
	"growlink/internal/domain/entity"
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

func (uc *AnswerReaderUsecase) GetAnswer(answerID uint) (*entity.Answer, error) {
	answer, err := uc.answerRepo.GetByID(answerID)
	if err != nil {
		return nil, fmt.Errorf("回答情報の取得に失敗しました: %w", err)
	}
	return answer, nil
}

func (uc *AnswerReaderUsecase) ListAnswersByMeeting(meetingID uint) ([]entity.Answer, error) {
	answers, err := uc.answerRepo.GetByMeetingID(meetingID)
	if err != nil {
		return nil, fmt.Errorf("回答情報の取得に失敗しました: %w", err)
	}
	return answers, nil
}
