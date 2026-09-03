package meeting_mng_usecase_test

import (
	"testing"
	"time"

	meeting_mng_usecase "growlink/internal/application/usecase/management_usecase/meeting_management_usecase"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
	repomock "growlink/internal/domain/interface/repository/mock"
	"growlink/internal/shared/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type meetingCreatorCaptures struct {
	createdMeetings []entity.Meeting
}

func newMeetingCreator(t *testing.T) (
	*meeting_mng_usecase.MeetingCreator,
	*repomock.MockMeetingRepository,
	*repomock.MockPairingRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	meetingRepo := repomock.NewMockMeetingRepository(ctrl)
	pairingRepo := repomock.NewMockPairingRepository(ctrl)
	return meeting_mng_usecase.NewMeetingCreator(meetingRepo, pairingRepo), meetingRepo, pairingRepo
}

func TestMeetingCreator_CreateMeeting(t *testing.T) {
	weekday := time.Date(2026, time.January, 2, 9, 0, 0, 0, time.UTC)
	holiday := time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC)
	expectedAfterHoliday := time.Date(2026, time.January, 2, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		pairingID     uint
		templateID    uint
		implementDate time.Time
		wantErr       bool
		errContains   string
		setup         func(meeting *repomock.MockMeetingRepository, cap *meetingCreatorCaptures)
		assertFn      func(t *testing.T, got *outputmodel.MeetingOutput, cap *meetingCreatorCaptures)
	}{
		{
			name:          "通常系_平日の実施日",
			pairingID:     pairingID1,
			templateID:    validTemplateID,
			implementDate: weekday,
			wantErr:       false,
			setup: func(meeting *repomock.MockMeetingRepository, cap *meetingCreatorCaptures) {
				meeting.EXPECT().Create(gomock.Any()).DoAndReturn(func(m *entity.Meeting) (*entity.Meeting, error) {
					cap.createdMeetings = append(cap.createdMeetings, *m)
					created := *m
					created.ID = meetingID1
					return &created, nil
				})
			},
			assertFn: func(t *testing.T, got *outputmodel.MeetingOutput, cap *meetingCreatorCaptures) {
				t.Helper()
				require.NotNil(t, got)
				require.Len(t, cap.createdMeetings, 1)
				assert.Equal(t, pairingID1, cap.createdMeetings[0].PairingID)
				assert.Equal(t, validTemplateID, cap.createdMeetings[0].TemplateID)
				assert.Equal(t, valueobject.MeetingStatusPending, cap.createdMeetings[0].Status)
				assert.True(t, cap.createdMeetings[0].ImplementedAt.Equal(weekday))
				assert.Equal(t, meetingID1, got.ID)
			},
		},
		{
			name:          "通常系_祝日は翌平日に繰り越す",
			pairingID:     pairingID1,
			templateID:    validTemplateID,
			implementDate: holiday,
			wantErr:       false,
			setup: func(meeting *repomock.MockMeetingRepository, cap *meetingCreatorCaptures) {
				meeting.EXPECT().Create(gomock.Any()).DoAndReturn(func(m *entity.Meeting) (*entity.Meeting, error) {
					cap.createdMeetings = append(cap.createdMeetings, *m)
					created := *m
					created.ID = meetingID1
					return &created, nil
				})
			},
			assertFn: func(t *testing.T, got *outputmodel.MeetingOutput, cap *meetingCreatorCaptures) {
				t.Helper()
				require.Len(t, cap.createdMeetings, 1)
				assert.True(t, cap.createdMeetings[0].ImplementedAt.Equal(expectedAfterHoliday))
			},
		},
		{
			name:          "異常系_Create失敗",
			pairingID:     pairingID1,
			templateID:    validTemplateID,
			implementDate: weekday,
			wantErr:       true,
			errContains:   "ミーティングの作成に失敗しました",
			setup: func(meeting *repomock.MockMeetingRepository, _ *meetingCreatorCaptures) {
				meeting.EXPECT().Create(gomock.Any()).Return(nil, errCreateMeetingFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, meetingRepo, _ := newMeetingCreator(t)
			cap := &meetingCreatorCaptures{}
			if tt.setup != nil {
				tt.setup(meetingRepo, cap)
			}

			got, gotErr := uc.CreateMeeting(tt.pairingID, tt.templateID, tt.implementDate)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "CreateMeeting")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, got, cap)
			}
		})
	}
}

func TestMeetingCreator_CreateMeetingEveryMonth(t *testing.T) {
	baseDate := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	expectedImplementDate := utils.GetSpecificWeekday(baseDate, 2, time.Wednesday)

	tests := []struct {
		name       string
		templateID uint
		baseDate   time.Time
		week       int
		weekday    time.Weekday
		wantErr    bool
		errContains string
		setup      func(meeting *repomock.MockMeetingRepository, pairing *repomock.MockPairingRepository, cap *meetingCreatorCaptures)
		assertFn   func(t *testing.T, got []outputmodel.MeetingOutput, cap *meetingCreatorCaptures)
	}{
		{
			name:       "通常系_全ペアリング分の定期ミーティング作成",
			templateID: validTemplateID,
			baseDate:   baseDate,
			week:       2,
			weekday:    time.Wednesday,
			wantErr:    false,
			setup: func(meeting *repomock.MockMeetingRepository, pairing *repomock.MockPairingRepository, cap *meetingCreatorCaptures) {
				pairing.EXPECT().GetAll().Return([]entity.Pairing{
					*newTestPairing(pairingID1, validMenteeID, validMentorID),
					*newTestPairing(pairingID2, validMemberID, validMentorID),
				}, nil)
				meeting.EXPECT().Create(gomock.Any()).DoAndReturn(func(m *entity.Meeting) (*entity.Meeting, error) {
					cap.createdMeetings = append(cap.createdMeetings, *m)
					created := *m
					if m.PairingID == pairingID1 {
						created.ID = meetingID1
					} else {
						created.ID = meetingID2
					}
					return &created, nil
				}).Times(2)
			},
			assertFn: func(t *testing.T, got []outputmodel.MeetingOutput, cap *meetingCreatorCaptures) {
				assertMeetingOutputIDSet(t, []uint{meetingID1, meetingID2}, got)
				require.Len(t, cap.createdMeetings, 2)
				for _, m := range cap.createdMeetings {
					assert.Equal(t, validTemplateID, m.TemplateID)
					assert.Equal(t, valueobject.MeetingStatusPending, m.Status)
					assert.True(t, m.ImplementedAt.Equal(expectedImplementDate))
				}
				pairingIDs := map[uint]struct{}{cap.createdMeetings[0].PairingID: {}, cap.createdMeetings[1].PairingID: {}}
				_, ok1 := pairingIDs[pairingID1]
				_, ok2 := pairingIDs[pairingID2]
				assert.True(t, ok1 && ok2)
			},
		},
		{
			name:       "通常系_ペアリングが空",
			templateID: validTemplateID,
			baseDate:   baseDate,
			week:       2,
			weekday:    time.Wednesday,
			wantErr:    false,
			setup: func(_ *repomock.MockMeetingRepository, pairing *repomock.MockPairingRepository, _ *meetingCreatorCaptures) {
				pairing.EXPECT().GetAll().Return([]entity.Pairing{}, nil)
			},
			assertFn: func(t *testing.T, got []outputmodel.MeetingOutput, _ *meetingCreatorCaptures) {
				assertMeetingOutputIDSet(t, nil, got)
			},
		},
		{
			name:        "異常系_GetAll失敗",
			templateID:  validTemplateID,
			baseDate:    baseDate,
			week:        2,
			weekday:     time.Wednesday,
			wantErr:     true,
			errContains: "定期ミーティング作成処理時に全てのペアリング情報の取得に失敗しました",
			setup: func(_ *repomock.MockMeetingRepository, pairing *repomock.MockPairingRepository, _ *meetingCreatorCaptures) {
				pairing.EXPECT().GetAll().Return(nil, errGetPairingsFailed)
			},
		},
		{
			name:        "異常系_Create失敗",
			templateID:  validTemplateID,
			baseDate:    baseDate,
			week:        2,
			weekday:     time.Wednesday,
			wantErr:     true,
			errContains: "定期ミーティングのペアリングのミーティング作成に失敗しました",
			setup: func(meeting *repomock.MockMeetingRepository, pairing *repomock.MockPairingRepository, _ *meetingCreatorCaptures) {
				pairing.EXPECT().GetAll().Return([]entity.Pairing{
					*newTestPairing(pairingID1, validMenteeID, validMentorID),
				}, nil)
				meeting.EXPECT().Create(gomock.Any()).Return(nil, errCreateMeetingFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, meetingRepo, pairingRepo := newMeetingCreator(t)
			cap := &meetingCreatorCaptures{}
			if tt.setup != nil {
				tt.setup(meetingRepo, pairingRepo, cap)
			}

			got, gotErr := uc.CreateMeetingEveryMonth(tt.templateID, tt.baseDate, tt.week, tt.weekday)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "CreateMeetingEveryMonth")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, got, cap)
			}
		})
	}
}

func TestMeetingCreator_CreateMeetingByAll(t *testing.T) {
	implementDate := time.Date(2026, time.January, 2, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		templateID    uint
		implementDate time.Time
		wantErr       bool
		errContains   string
		setup         func(meeting *repomock.MockMeetingRepository, pairing *repomock.MockPairingRepository, cap *meetingCreatorCaptures)
		assertFn      func(t *testing.T, got []outputmodel.MeetingOutput, cap *meetingCreatorCaptures)
	}{
		{
			name:          "通常系_全ペアリング分作成",
			templateID:    validTemplateID,
			implementDate: implementDate,
			wantErr:       false,
			setup: func(meeting *repomock.MockMeetingRepository, pairing *repomock.MockPairingRepository, cap *meetingCreatorCaptures) {
				pairing.EXPECT().GetAll().Return([]entity.Pairing{
					*newTestPairing(pairingID1, validMenteeID, validMentorID),
					*newTestPairing(pairingID2, validMemberID, validMentorID),
				}, nil)
				meeting.EXPECT().Create(gomock.Any()).DoAndReturn(func(m *entity.Meeting) (*entity.Meeting, error) {
					cap.createdMeetings = append(cap.createdMeetings, *m)
					created := *m
					if m.PairingID == pairingID1 {
						created.ID = meetingID1
					} else {
						created.ID = meetingID2
					}
					return &created, nil
				}).Times(2)
			},
			assertFn: func(t *testing.T, got []outputmodel.MeetingOutput, cap *meetingCreatorCaptures) {
				assertMeetingOutputIDSet(t, []uint{meetingID1, meetingID2}, got)
				require.Len(t, cap.createdMeetings, 2)
				pairingIDs := map[uint]struct{}{cap.createdMeetings[0].PairingID: {}, cap.createdMeetings[1].PairingID: {}}
				_, ok1 := pairingIDs[pairingID1]
				_, ok2 := pairingIDs[pairingID2]
				assert.True(t, ok1 && ok2)
			},
		},
		{
			name:          "通常系_ペアリングが空",
			templateID:    validTemplateID,
			implementDate: implementDate,
			wantErr:       false,
			setup: func(_ *repomock.MockMeetingRepository, pairing *repomock.MockPairingRepository, _ *meetingCreatorCaptures) {
				pairing.EXPECT().GetAll().Return([]entity.Pairing{}, nil)
			},
			assertFn: func(t *testing.T, got []outputmodel.MeetingOutput, _ *meetingCreatorCaptures) {
				assertMeetingOutputIDSet(t, nil, got)
			},
		},
		{
			name:          "異常系_GetAll失敗",
			templateID:    validTemplateID,
			implementDate: implementDate,
			wantErr:       true,
			errContains:   "ミーティング作成時に全てのペアリング情報の取得に失敗しました",
			setup: func(_ *repomock.MockMeetingRepository, pairing *repomock.MockPairingRepository, _ *meetingCreatorCaptures) {
				pairing.EXPECT().GetAll().Return(nil, errGetPairingsFailed)
			},
		},
		{
			name:          "異常系_Create失敗",
			templateID:    validTemplateID,
			implementDate: implementDate,
			wantErr:       true,
			errContains:   "全ペアリングのミーティング作成に失敗しました",
			setup: func(meeting *repomock.MockMeetingRepository, pairing *repomock.MockPairingRepository, _ *meetingCreatorCaptures) {
				pairing.EXPECT().GetAll().Return([]entity.Pairing{
					*newTestPairing(pairingID1, validMenteeID, validMentorID),
				}, nil)
				meeting.EXPECT().Create(gomock.Any()).Return(nil, errCreateMeetingFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, meetingRepo, pairingRepo := newMeetingCreator(t)
			cap := &meetingCreatorCaptures{}
			if tt.setup != nil {
				tt.setup(meetingRepo, pairingRepo, cap)
			}

			got, gotErr := uc.CreateMeetingByAll(tt.templateID, tt.implementDate)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "CreateMeetingByAll")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, got, cap)
			}
		})
	}
}
