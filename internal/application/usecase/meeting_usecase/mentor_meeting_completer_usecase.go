package meeting_usecase

import (
	"fmt"

	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
)

type MentorTopicCommentInput struct {
	AnswerUUID string `json:"answerUuid"`
	Content    string `json:"content"`
}

type MentorMeetingCompleterUsecase struct {
	meetingRepo       repository_interface.MeetingRepository
	answerRepo        repository_interface.AnswerRepository
	answerCommenter   *AnswerCommenterUsecase
}

func NewMentorMeetingCompleterUsecase(
	meetingRepo repository_interface.MeetingRepository,
	answerRepo repository_interface.AnswerRepository,
	answerCommenter *AnswerCommenterUsecase,
) *MentorMeetingCompleterUsecase {
	return &MentorMeetingCompleterUsecase{
		meetingRepo:     meetingRepo,
		answerRepo:      answerRepo,
		answerCommenter: answerCommenter,
	}
}

func (uc *MentorMeetingCompleterUsecase) Complete(
	meetingID, mentorUserID, menteeUserID uint,
	comments []MentorTopicCommentInput,
) error {
	for _, item := range comments {
		if item.AnswerUUID == "" {
			continue
		}
		answerUUID, err := uuid.Parse(item.AnswerUUID)
		if err != nil {
			return fmt.Errorf("回答UUIDが不正です: %w", err)
		}
		answer, err := uc.answerRepo.GetByUUID(answerUUID)
		if err != nil {
			return fmt.Errorf("回答の取得に失敗しました: %w", err)
		}
		if answer == nil {
			return fmt.Errorf("回答が見つかりません")
		}
		if answer.MeetingID != meetingID {
			return fmt.Errorf("この面談の回答ではありません")
		}
		if answer.UserID != menteeUserID {
			return fmt.Errorf("メンティーの回答ではありません")
		}
		if err := uc.answerCommenter.SaveComment(answer.ID, mentorUserID, item.Content); err != nil {
			return err
		}
	}
	if err := uc.meetingRepo.UpdateStatus(meetingID, valueobject.MeetingStatusDone); err != nil {
		return fmt.Errorf("面談の完了更新に失敗しました: %w", err)
	}
	return nil
}
