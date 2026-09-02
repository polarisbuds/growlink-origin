package usecase_interface

import "growlink/internal/domain/entity"

type MeetingUsecase interface {
	AnswerReader
	DraftSaver
	SubmitCanceller
	Submitter

	AnswerCommenter
	MeetingCommenter
}

type AnswerReader interface {
	GetAnswer(answerID uint) (*entity.Answer, error)
	ListAnswersByMeeting(meetingID uint) ([]entity.Answer, error)
}

type DraftSaver interface {
	SaveDraft(meetingID uint, answers []entity.Answer) error
}

type Submitter interface {
	Submit(meetingID uint, answers []entity.Answer) error
}

type SubmitCanceller interface {
	CancelSubmit(meetingID uint) (bool, error)
}

type AnswerCommenter interface {
	PostComment(answerID uint, content string) error
	EditComment(commentID uint, content string) error
	DeleteComment(commentID uint) error
}

type MeetingCommenter interface {
	PostComment(meetingID uint, content string) error
	EditComment(commentID uint, content string) error
	DeleteComment(commentID uint) error
}
