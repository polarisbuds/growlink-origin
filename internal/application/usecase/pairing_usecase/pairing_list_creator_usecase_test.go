package pairing_usecase_test

import (
	"testing"

	"growlink/internal/application/usecase/pairing_usecase"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type pairingListCreatorSetupFn func(t *testing.T, pairing *repomock.MockPairingRepository)
type pairingListCreatorAssertFn func(t *testing.T, got []outputmodel.UserOutput)

func newPairingListCreatorUsecase(t *testing.T) (
	*pairing_usecase.PairingListCreatorUsecase,
	*repomock.MockPairingRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	pairingRepo := repomock.NewMockPairingRepository(ctrl)
	userRepo := repomock.NewMockUserRepository(ctrl)
	return pairing_usecase.NewPairingListCreatorUsecase(pairingRepo, userRepo), pairingRepo
}

func TestPairingListCreatorUsecase_CreatePairList(t *testing.T) {
	menteeA := newTestUser(31, "mentee-a")
	menteeB := newTestUser(32, "mentee-b")
	pairingA := newTestPairing(301, menteeA.ID, validMentorID)
	pairingA.Mentee = menteeA
	pairingB := newTestPairing(302, menteeB.ID, validMentorID)
	pairingB.Mentee = menteeB

	tests := []struct {
		name        string
		mentorID    uint
		wantErr     bool
		errContains string
		setup       pairingListCreatorSetupFn
		assertFn    pairingListCreatorAssertFn
	}{
		{
			name:     "通常系_メンティー一覧",
			mentorID: validMentorID,
			wantErr:  false,
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository) {
				pairing.EXPECT().GetByMentorID(validMentorID).Return([]entity.Pairing{*pairingA, *pairingB}, nil)
			},
			assertFn: func(t *testing.T, got []outputmodel.UserOutput) {
				t.Helper()
				require.Len(t, got, 2)
				assert.Equal(t, menteeA.ID, got[0].ID)
				assert.Equal(t, menteeA.Name, got[0].Name)
				assert.Equal(t, menteeB.ID, got[1].ID)
				assert.Equal(t, menteeB.Name, got[1].Name)
			},
		},
		{
			name:     "通常系_0件",
			mentorID: validMentorID,
			wantErr:  false,
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository) {
				pairing.EXPECT().GetByMentorID(validMentorID).Return([]entity.Pairing{}, nil)
			},
			assertFn: func(t *testing.T, got []outputmodel.UserOutput) {
				t.Helper()
				assert.Empty(t, got)
			},
		},
		{
			name:        "異常系_取得失敗",
			mentorID:    validMentorID,
			wantErr:     true,
			errContains: "ペアリストの作成に失敗しました",
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository) {
				pairing.EXPECT().GetByMentorID(validMentorID).Return(nil, errGetPairingFailed)
			},
			assertFn: func(t *testing.T, got []outputmodel.UserOutput) {
				t.Helper()
				assert.Nil(t, got)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, pairingRepo := newPairingListCreatorUsecase(t)
			if tt.setup != nil {
				tt.setup(t, pairingRepo)
			}

			got, gotErr := uc.CreatePairList(tt.mentorID)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "CreatePairList")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, got)
			}
		})
	}
}
