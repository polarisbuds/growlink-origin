package usecase_interface

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import "growlink/internal/domain/entity"

type (
	MeetingUsecase interface {
		AnswerReader
		DraftSaver
		SubmitCanceller
		Submitter

		AnswerCommenter
		MeetingCommenter
	}

	AnswerReader interface {
		GetAnswer(answerID uint) (*entity.Answer, error)
		ListAnswersByMeeting(meetingID uint) ([]entity.Answer, error)
	}
	DraftSaver interface {
		SaveDraft(meetingID uint, answers []entity.Answer) error
	}
	Submitter interface {
		Submit(meetingID uint, answers []entity.Answer) error
	}
	SubmitCanceller interface {
		CancelSubmit(meetingID uint) error
	}
	AnswerCommenter interface {
		PostComment(answerID uint, content string) error
		EditComment(commentID uint, content string) error
		DeleteComment(commentID uint) error
	}
	MeetingCommenter interface {
		PostComment(meetingID uint, content string) error
		EditComment(commentID uint, content string) error
		DeleteComment(commentID uint) error
	}
)
