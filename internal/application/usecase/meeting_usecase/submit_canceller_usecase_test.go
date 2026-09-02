package meeting_usecase_test

import (
	"testing"

	"growlink/internal/application/usecase/meeting_usecase"
	repomock "growlink/internal/domain/interface/repository/mock"
	"growlink/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type submitCancellerCaptures struct {
	updatedMeetingID uint
	updatedStatus    valueobject.MeetingStatus
}

type submitCancellerSetupFn func(t *testing.T, meeting *repomock.MockMeetingRepository, cap *submitCancellerCaptures)
type submitCancellerAssertFn func(t *testing.T, cap *submitCancellerCaptures)

func newSubmitCancellerUsecase(t *testing.T) (
	*meeting_usecase.SubmitCancellerUsecase,
	*repomock.MockMeetingRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	meetingRepo := repomock.NewMockMeetingRepository(ctrl)
	return meeting_usecase.NewSubmitCancellerUsecase(meetingRepo), meetingRepo
}

func TestSubmitCancellerUsecase_CancelSubmit(t *testing.T) {
	tests := []struct {
		name        string
		meetingID   uint
		wantErr     bool
		errContains string
		setup       submitCancellerSetupFn
		assertFn    submitCancellerAssertFn
	}{
		{
			name:      "通常系_提出取消",
			meetingID: validMeetingID,
			wantErr:   false,
			setup: func(_ *testing.T, meeting *repomock.MockMeetingRepository, cap *submitCancellerCaptures) {
				meeting.EXPECT().UpdateStatus(validMeetingID, valueobject.MeetingStatusDraft).DoAndReturn(
					func(meetingID uint, status valueobject.MeetingStatus) error {
						cap.updatedMeetingID = meetingID
						cap.updatedStatus = status
						return nil
					},
				)
			},
			assertFn: func(t *testing.T, cap *submitCancellerCaptures) {
				t.Helper()
				assert.Equal(t, validMeetingID, cap.updatedMeetingID)
				assert.Equal(t, valueobject.MeetingStatusDraft, cap.updatedStatus)
			},
		},
		{
			name:        "異常系_ステータス更新失敗",
			meetingID:   validMeetingID,
			wantErr:     true,
			errContains: "下書き保存へのステータス更新に失敗しました",
			setup: func(_ *testing.T, meeting *repomock.MockMeetingRepository, _ *submitCancellerCaptures) {
				meeting.EXPECT().UpdateStatus(validMeetingID, valueobject.MeetingStatusDraft).Return(errUpdateStatusFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, meetingRepo := newSubmitCancellerUsecase(t)
			cap := &submitCancellerCaptures{}
			if tt.setup != nil {
				tt.setup(t, meetingRepo, cap)
			}

			err := uc.CancelSubmit(tt.meetingID)
			assertUsecaseErr(t, err, tt.wantErr, tt.errContains, "CancelSubmit")
			if !t.Failed() && tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}
