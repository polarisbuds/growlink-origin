package meeting_mng_usecase_test

import (
	"testing"
	"time"

	meeting_mng_usecase "growlink/internal/application/usecase/management_usecase/meeting_management_usecase"
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type meetingUpdaterCaptures struct {
	meetingID     uint
	meeting       *entity.Meeting
	implementDate time.Time
	pairingID     uint
	templateID    uint
}

func newMeetingUpdaterUsecase(t *testing.T) (
	*meeting_mng_usecase.MeetingUpdaterUsecase,
	*repomock.MockMeetingRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	meetingRepo := repomock.NewMockMeetingRepository(ctrl)
	return meeting_mng_usecase.NewMeetingUpdaterUsecase(meetingRepo), meetingRepo
}

func TestMeetingUpdaterUsecase_Update(t *testing.T) {
	meeting := &entity.Meeting{
		Status:     valueobject.MeetingStatusDone,
		TemplateID: validTemplateID,
		PairingID:  pairingID1,
	}

	tests := []struct {
		name        string
		meetingID   uint
		meeting     *entity.Meeting
		wantErr     bool
		errContains string
		setup       func(repo *repomock.MockMeetingRepository, cap *meetingUpdaterCaptures)
		assertFn    func(t *testing.T, cap *meetingUpdaterCaptures)
	}{
		{
			name:      "通常系_更新",
			meetingID: meetingID1,
			meeting:   meeting,
			wantErr:   false,
			setup: func(repo *repomock.MockMeetingRepository, cap *meetingUpdaterCaptures) {
				repo.EXPECT().Update(meetingID1, meeting).DoAndReturn(func(id uint, m *entity.Meeting) error {
					cap.meetingID = id
					cap.meeting = m
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *meetingUpdaterCaptures) {
				assert.Equal(t, meetingID1, cap.meetingID)
				assert.Equal(t, valueobject.MeetingStatusDone, cap.meeting.Status)
			},
		},
		{
			name:        "異常系_更新失敗",
			meetingID:   meetingID1,
			meeting:     meeting,
			wantErr:     true,
			errContains: "ミーティングの更新に失敗しました",
			setup: func(repo *repomock.MockMeetingRepository, _ *meetingUpdaterCaptures) {
				repo.EXPECT().Update(meetingID1, meeting).Return(errUpdateMeetingFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo := newMeetingUpdaterUsecase(t)
			cap := &meetingUpdaterCaptures{}
			if tt.setup != nil {
				tt.setup(repo, cap)
			}

			gotErr := uc.Update(tt.meetingID, tt.meeting)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "Update")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}

func TestMeetingUpdaterUsecase_ChangeImplementDate(t *testing.T) {
	newDate := time.Date(2026, time.March, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		meetingID     uint
		implementDate time.Time
		wantErr       bool
		errContains   string
		setup         func(repo *repomock.MockMeetingRepository, cap *meetingUpdaterCaptures)
		assertFn      func(t *testing.T, cap *meetingUpdaterCaptures)
	}{
		{
			name:          "通常系_実施日変更",
			meetingID:     meetingID1,
			implementDate: newDate,
			wantErr:       false,
			setup: func(repo *repomock.MockMeetingRepository, cap *meetingUpdaterCaptures) {
				repo.EXPECT().UpdateImplementDate(meetingID1, newDate).DoAndReturn(func(id uint, d time.Time) error {
					cap.meetingID = id
					cap.implementDate = d
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *meetingUpdaterCaptures) {
				assert.Equal(t, meetingID1, cap.meetingID)
				assert.True(t, cap.implementDate.Equal(newDate))
			},
		},
		{
			name:          "異常系_更新失敗",
			meetingID:     meetingID1,
			implementDate: newDate,
			wantErr:       true,
			errContains:   "ミーティングの実施予定日の変更に失敗しました",
			setup: func(repo *repomock.MockMeetingRepository, _ *meetingUpdaterCaptures) {
				repo.EXPECT().UpdateImplementDate(meetingID1, newDate).Return(errUpdateMeetingFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo := newMeetingUpdaterUsecase(t)
			cap := &meetingUpdaterCaptures{}
			if tt.setup != nil {
				tt.setup(repo, cap)
			}

			gotErr := uc.ChangeImplementDate(tt.meetingID, tt.implementDate)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "ChangeImplementDate")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}

func TestMeetingUpdaterUsecase_ChangePairing(t *testing.T) {
	newPairingID := pairingID2

	tests := []struct {
		name         string
		meetingID    uint
		newPairingID uint
		wantErr      bool
		errContains  string
		setup        func(repo *repomock.MockMeetingRepository, cap *meetingUpdaterCaptures)
		assertFn     func(t *testing.T, cap *meetingUpdaterCaptures)
	}{
		{
			name:         "通常系_ペアリング変更",
			meetingID:    meetingID1,
			newPairingID: newPairingID,
			wantErr:      false,
			setup: func(repo *repomock.MockMeetingRepository, cap *meetingUpdaterCaptures) {
				repo.EXPECT().UpdatePairing(meetingID1, newPairingID).DoAndReturn(func(id, pid uint) error {
					cap.meetingID = id
					cap.pairingID = pid
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *meetingUpdaterCaptures) {
				assert.Equal(t, meetingID1, cap.meetingID)
				assert.Equal(t, newPairingID, cap.pairingID)
			},
		},
		{
			name:         "異常系_更新失敗",
			meetingID:    meetingID1,
			newPairingID: newPairingID,
			wantErr:      true,
			errContains:  "ミーティングのペアリングの変更に失敗しました",
			setup: func(repo *repomock.MockMeetingRepository, _ *meetingUpdaterCaptures) {
				repo.EXPECT().UpdatePairing(meetingID1, newPairingID).Return(errUpdateMeetingFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo := newMeetingUpdaterUsecase(t)
			cap := &meetingUpdaterCaptures{}
			if tt.setup != nil {
				tt.setup(repo, cap)
			}

			gotErr := uc.ChangePairing(tt.meetingID, tt.newPairingID)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "ChangePairing")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}

func TestMeetingUpdaterUsecase_ChangeTemplate(t *testing.T) {
	newTemplateID := uint(999)

	tests := []struct {
		name          string
		meetingID     uint
		newTemplateID uint
		wantErr       bool
		errContains   string
		setup         func(repo *repomock.MockMeetingRepository, cap *meetingUpdaterCaptures)
		assertFn      func(t *testing.T, cap *meetingUpdaterCaptures)
	}{
		{
			name:          "通常系_テンプレート変更",
			meetingID:     meetingID1,
			newTemplateID: newTemplateID,
			wantErr:       false,
			setup: func(repo *repomock.MockMeetingRepository, cap *meetingUpdaterCaptures) {
				repo.EXPECT().UpdateTemplate(meetingID1, newTemplateID).DoAndReturn(func(id, tid uint) error {
					cap.meetingID = id
					cap.templateID = tid
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *meetingUpdaterCaptures) {
				assert.Equal(t, meetingID1, cap.meetingID)
				assert.Equal(t, newTemplateID, cap.templateID)
			},
		},
		{
			name:          "異常系_更新失敗",
			meetingID:     meetingID1,
			newTemplateID: newTemplateID,
			wantErr:       true,
			errContains:   "ミーティングのテンプレートの変更に失敗しました",
			setup: func(repo *repomock.MockMeetingRepository, _ *meetingUpdaterCaptures) {
				repo.EXPECT().UpdateTemplate(meetingID1, newTemplateID).Return(errUpdateMeetingFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo := newMeetingUpdaterUsecase(t)
			cap := &meetingUpdaterCaptures{}
			if tt.setup != nil {
				tt.setup(repo, cap)
			}

			gotErr := uc.ChangeTemplate(tt.meetingID, tt.newTemplateID)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "ChangeTemplate")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}
