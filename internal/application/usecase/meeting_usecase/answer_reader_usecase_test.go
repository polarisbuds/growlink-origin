package meeting_usecase_test

import (
	"testing"

	"growlink/internal/application/usecase/meeting_usecase"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type answerReaderSetupFn func(t *testing.T, answer *repomock.MockAnswerRepository)
type answerReaderAssertFn func(t *testing.T, got *entity.Answer, gotList []entity.Answer)

func newAnswerReaderUsecase(t *testing.T) (
	*meeting_usecase.AnswerReaderUsecase,
	*repomock.MockAnswerRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	answerRepo := repomock.NewMockAnswerRepository(ctrl)
	return meeting_usecase.NewAnswerReaderUsecase(answerRepo), answerRepo
}

func TestAnswerReaderUsecase_GetAnswer(t *testing.T) {
	tests := []struct {
		name        string
		answerID    uint
		wantErr     bool
		errContains string
		setup       answerReaderSetupFn
		assertFn    answerReaderAssertFn
	}{
		{
			name:     "通常系_回答取得",
			answerID: validAnswerID,
			wantErr:  false,
			setup: func(_ *testing.T, answer *repomock.MockAnswerRepository) {
				answer.EXPECT().GetByID(validAnswerID).Return(newTestAnswer(validAnswerID, validMeetingID, "answer"), nil)
			},
			assertFn: func(t *testing.T, got *entity.Answer, _ []entity.Answer) {
				t.Helper()
				require.NotNil(t, got)
				assert.Equal(t, validAnswerID, got.ID)
				assert.Equal(t, "answer", got.Content)
			},
		},
		{
			name:        "異常系_リポジトリエラー",
			answerID:    validAnswerID,
			wantErr:     true,
			errContains: "回答情報の取得に失敗しました",
			setup: func(_ *testing.T, answer *repomock.MockAnswerRepository) {
				answer.EXPECT().GetByID(validAnswerID).Return(nil, errGetAnswerFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, answerRepo := newAnswerReaderUsecase(t)
			if tt.setup != nil {
				tt.setup(t, answerRepo)
			}

			got, err := uc.GetAnswer(tt.answerID)
			assertUsecaseErr(t, err, tt.wantErr, tt.errContains, "GetAnswer")
			if !t.Failed() && tt.assertFn != nil {
				tt.assertFn(t, got, nil)
			}
		})
	}
}

func TestAnswerReaderUsecase_ListAnswersByMeeting(t *testing.T) {
	tests := []struct {
		name        string
		meetingID   uint
		wantErr     bool
		errContains string
		setup       answerReaderSetupFn
		assertFn    answerReaderAssertFn
	}{
		{
			name:      "通常系_ミーティング別回答一覧",
			meetingID: validMeetingID,
			wantErr:   false,
			setup: func(_ *testing.T, answer *repomock.MockAnswerRepository) {
				answers := []entity.Answer{
					*newTestAnswer(1, validMeetingID, "a1"),
					*newTestAnswer(2, validMeetingID, "a2"),
				}
				answer.EXPECT().GetByMeetingID(validMeetingID).Return(answers, nil)
			},
			assertFn: func(t *testing.T, _ *entity.Answer, gotList []entity.Answer) {
				t.Helper()
				require.Len(t, gotList, 2)
				assert.Equal(t, "a1", gotList[0].Content)
				assert.Equal(t, "a2", gotList[1].Content)
			},
		},
		{
			name:        "異常系_リポジトリエラー",
			meetingID:   validMeetingID,
			wantErr:     true,
			errContains: "回答情報の取得に失敗しました",
			setup: func(_ *testing.T, answer *repomock.MockAnswerRepository) {
				answer.EXPECT().GetByMeetingID(validMeetingID).Return(nil, errGetByMeetingFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, answerRepo := newAnswerReaderUsecase(t)
			if tt.setup != nil {
				tt.setup(t, answerRepo)
			}

			gotList, err := uc.ListAnswersByMeeting(tt.meetingID)
			assertUsecaseErr(t, err, tt.wantErr, tt.errContains, "ListAnswersByMeeting")
			if !t.Failed() && tt.assertFn != nil {
				tt.assertFn(t, nil, gotList)
			}
		})
	}
}
