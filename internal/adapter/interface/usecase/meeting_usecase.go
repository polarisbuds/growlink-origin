package usecase_interface

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import (
	outputmodel "growlink/internal/application/model/output_model"
)

type (
	MeetingUsecase interface {
		AnswerReader
		DraftSaver
		SubmitCanceller
		Submitter

		AnswerCommenter
		MeetingCommenter
	}
	MeetingList interface {
		ListMyMeeting(menteeID uint) ([]outputmodel.MeetingOutput, error)
	}
	AnswerReader interface {
		GetAnswer(answerID uint) (*outputmodel.AnswerOutput, error)
		ListAnswersByMeeting(meetingID uint) ([]outputmodel.AnswerOutput, error)
	}
	DraftSaver interface {
		SaveDraft(meetingID uint, answers []outputmodel.AnswerOutput) error
	}
	Submitter interface {
		Submit(meetingID uint, answers []outputmodel.AnswerOutput) error
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
