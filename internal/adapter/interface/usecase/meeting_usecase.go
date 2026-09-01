package usecase_interface

import "growlink/internal/domain/entity"

type MeetingUsecase interface {
	AnswerReader
	DraftSaver
	Submitter
}

type AnswerReader interface {
	GetAnswer() (*entity.Answer, error)
	ListAnswersByMeeting(meetingID uint) ([]entity.Answer, error)
}

type DraftSaver interface {
	SaveDraft(answers []entity.Answer) error
}

type Submitter interface {
	Submit(answers []entity.Answer) error
}
