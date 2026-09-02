package meeting_usecase_test

import (
	"testing"

	"growlink/internal/application/usecase/meeting_usecase"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"
	"growlink/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type submitterCaptures struct {
	submittedAnswers []entity.Answer
	updatedMeetingID uint
	updatedStatus    valueobject.MeetingStatus
}

type submitterSetupFn func(t *testing.T, meeting *repomock.MockMeetingRepository, answer *repomock.MockAnswerRepository, cap *submitterCaptures)
type submitterAssertFn func(t *testing.T, cap *submitterCaptures)

func newSubmitterUsecase(t *testing.T) (
	*meeting_usecase.SubmitterUsecase,
	*repomock.MockMeetingRepository,
	*repomock.MockAnswerRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	meetingRepo := repomock.NewMockMeetingRepository(ctrl)
	answerRepo := repomock.NewMockAnswerRepository(ctrl)
	return meeting_usecase.NewSubmitterUsecase(meetingRepo, answerRepo), meetingRepo, answerRepo
}

func TestSubmitterUsecase_Submit(t *testing.T) {
	answers := []entity.Answer{*newTestAnswer(1, validMeetingID, "submit")}

	tests := []struct {
		name        string
		meetingID   uint
		answers     []entity.Answer
		wantErr     bool
		errContains string
		setup       submitterSetupFn
		assertFn    submitterAssertFn
	}{
		{
			name:      "通常系_投稿",
			meetingID: validMeetingID,
			answers:   answers,
			wantErr:   false,
			setup: func(_ *testing.T, meeting *repomock.MockMeetingRepository, answer *repomock.MockAnswerRepository, cap *submitterCaptures) {
				answer.EXPECT().Submits(answers).DoAndReturn(func(submitted []entity.Answer) ([]entity.Answer, error) {
					cap.submittedAnswers = submitted
					return submitted, nil
				})
				meeting.EXPECT().UpdateStatus(validMeetingID, valueobject.MeetingStatusDone).DoAndReturn(
					func(meetingID uint, status valueobject.MeetingStatus) error {
						cap.updatedMeetingID = meetingID
						cap.updatedStatus = status
						return nil
					},
				)
			},
			assertFn: func(t *testing.T, cap *submitterCaptures) {
				t.Helper()
				assert.Equal(t, answers, cap.submittedAnswers)
				assert.Equal(t, validMeetingID, cap.updatedMeetingID)
				assert.Equal(t, valueobject.MeetingStatusDone, cap.updatedStatus)
			},
		},
		{
			name:        "異常系_回答投稿失敗",
			meetingID:   validMeetingID,
			answers:     answers,
			wantErr:     true,
			errContains: "回答の投稿に失敗しました",
			setup: func(_ *testing.T, _ *repomock.MockMeetingRepository, answer *repomock.MockAnswerRepository, _ *submitterCaptures) {
				answer.EXPECT().Submits(answers).Return(nil, errSubmitsFailed)
			},
		},
		{
			name:        "異常系_ステータス更新失敗",
			meetingID:   validMeetingID,
			answers:     answers,
			wantErr:     true,
			errContains: "投稿完了へのステータス更新に失敗しました",
			setup: func(_ *testing.T, meeting *repomock.MockMeetingRepository, answer *repomock.MockAnswerRepository, _ *submitterCaptures) {
				answer.EXPECT().Submits(answers).Return(answers, nil)
				meeting.EXPECT().UpdateStatus(validMeetingID, valueobject.MeetingStatusDone).Return(errUpdateStatusFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, meetingRepo, answerRepo := newSubmitterUsecase(t)
			cap := &submitterCaptures{}
			if tt.setup != nil {
				tt.setup(t, meetingRepo, answerRepo, cap)
			}

			err := uc.Submit(tt.meetingID, tt.answers)
			assertUsecaseErr(t, err, tt.wantErr, tt.errContains, "Submit")
			if !t.Failed() && tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}
