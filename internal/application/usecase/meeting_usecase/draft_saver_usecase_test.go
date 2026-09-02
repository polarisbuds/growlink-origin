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

type draftSaverCaptures struct {
	submittedAnswers []entity.Answer
	updatedMeetingID uint
	updatedStatus    valueobject.MeetingStatus
}

type draftSaverSetupFn func(t *testing.T, meeting *repomock.MockMeetingRepository, answer *repomock.MockAnswerRepository, cap *draftSaverCaptures)
type draftSaverAssertFn func(t *testing.T, cap *draftSaverCaptures)

func newDraftSaverUsecase(t *testing.T) (
	*meeting_usecase.DraftSaverUsecase,
	*repomock.MockMeetingRepository,
	*repomock.MockAnswerRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	meetingRepo := repomock.NewMockMeetingRepository(ctrl)
	answerRepo := repomock.NewMockAnswerRepository(ctrl)
	return meeting_usecase.NewDraftSaverUsecase(meetingRepo, answerRepo), meetingRepo, answerRepo
}

func TestDraftSaverUsecase_SaveDraft(t *testing.T) {
	answers := []entity.Answer{*newTestAnswer(1, validMeetingID, "draft")}

	tests := []struct {
		name        string
		meetingID   uint
		answers     []entity.Answer
		wantErr     bool
		errContains string
		setup       draftSaverSetupFn
		assertFn    draftSaverAssertFn
	}{
		{
			name:      "通常系_下書き保存",
			meetingID: validMeetingID,
			answers:   answers,
			wantErr:   false,
			setup: func(_ *testing.T, meeting *repomock.MockMeetingRepository, answer *repomock.MockAnswerRepository, cap *draftSaverCaptures) {
				answer.EXPECT().Submits(answers).DoAndReturn(func(submitted []entity.Answer) ([]entity.Answer, error) {
					cap.submittedAnswers = submitted
					return submitted, nil
				})
				meeting.EXPECT().UpdateStatus(validMeetingID, valueobject.MeetingStatusDraft).DoAndReturn(
					func(meetingID uint, status valueobject.MeetingStatus) error {
						cap.updatedMeetingID = meetingID
						cap.updatedStatus = status
						return nil
					},
				)
			},
			assertFn: func(t *testing.T, cap *draftSaverCaptures) {
				t.Helper()
				assert.Equal(t, answers, cap.submittedAnswers)
				assert.Equal(t, validMeetingID, cap.updatedMeetingID)
				assert.Equal(t, valueobject.MeetingStatusDraft, cap.updatedStatus)
			},
		},
		{
			name:        "異常系_回答保存失敗",
			meetingID:   validMeetingID,
			answers:     answers,
			wantErr:     true,
			errContains: "回答の下書き保存に失敗しました",
			setup: func(_ *testing.T, _ *repomock.MockMeetingRepository, answer *repomock.MockAnswerRepository, _ *draftSaverCaptures) {
				answer.EXPECT().Submits(answers).Return(nil, errSubmitsFailed)
			},
		},
		{
			name:        "異常系_ステータス更新失敗",
			meetingID:   validMeetingID,
			answers:     answers,
			wantErr:     true,
			errContains: "下書き保存へのステータス更新に失敗しました",
			setup: func(_ *testing.T, meeting *repomock.MockMeetingRepository, answer *repomock.MockAnswerRepository, _ *draftSaverCaptures) {
				answer.EXPECT().Submits(answers).Return(answers, nil)
				meeting.EXPECT().UpdateStatus(validMeetingID, valueobject.MeetingStatusDraft).Return(errUpdateStatusFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, meetingRepo, answerRepo := newDraftSaverUsecase(t)
			cap := &draftSaverCaptures{}
			if tt.setup != nil {
				tt.setup(t, meetingRepo, answerRepo, cap)
			}

			err := uc.SaveDraft(tt.meetingID, tt.answers)
			assertUsecaseErr(t, err, tt.wantErr, tt.errContains, "SaveDraft")
			if !t.Failed() && tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}
