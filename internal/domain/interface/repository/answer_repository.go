package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type AnswerRepository interface {
	GetByID(id uint) (*entity.Answer, error)
	GetByMeetingID(meetingID uint) ([]entity.Answer, error)
	GetByUserID(userD uint) ([]entity.Answer, error)

	Submit(answer *entity.Answer) (*entity.Answer, error)
	Submits(answers []entity.Answer) ([]entity.Answer, error)
}
