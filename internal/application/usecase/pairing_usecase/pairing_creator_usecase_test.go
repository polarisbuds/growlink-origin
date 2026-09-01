package pairing_usecase_test

import (
	"testing"

	"growlink/internal/application/usecase/pairing_usecase"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type pairingCreatorCaptures struct {
	createMenteeID uint
	createMentorID uint
	breakPairingID uint
	created        *entity.Pairing
}

type pairingCreatorSetupFn func(t *testing.T, pairing *repomock.MockPairingRepository, cap *pairingCreatorCaptures)
type pairingCreatorAssertFn func(t *testing.T, cap *pairingCreatorCaptures)

func newPairingCreatorUsecase(t *testing.T) (
	*pairing_usecase.PairingCreatorUsecase,
	*repomock.MockPairingRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	pairingRepo := repomock.NewMockPairingRepository(ctrl)
	userRepo := repomock.NewMockUserRepository(ctrl)
	return pairing_usecase.NewPairingCreatorUsecase(pairingRepo, userRepo), pairingRepo
}

func TestPairingCreatorUsecase_PairingCreator(t *testing.T) {
	tests := []struct {
		name        string
		menteeID    uint
		mentorID    uint
		wantErr     bool
		errContains string
		setup       pairingCreatorSetupFn
		assertFn    pairingCreatorAssertFn
	}{
		{
			name:     "通常系_旧ペアなし",
			menteeID: validMenteeID,
			mentorID: validMentorID,
			wantErr:  false,
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository, cap *pairingCreatorCaptures) {
				pairing.EXPECT().GetByMenteeID(validMenteeID).Return(nil, nil)
				pairing.EXPECT().CreatePairing(validMenteeID, validMentorID).DoAndReturn(
					func(menteeID, mentorID uint) (*entity.Pairing, error) {
						cap.createMenteeID = menteeID
						cap.createMentorID = mentorID
						created := newTestPairing(201, menteeID, mentorID)
						cap.created = created
						return created, nil
					},
				)
			},
			assertFn: func(t *testing.T, cap *pairingCreatorCaptures) {
				t.Helper()
				assert.Equal(t, validMenteeID, cap.createMenteeID)
				assert.Equal(t, validMentorID, cap.createMentorID)
				assert.Equal(t, uint(0), cap.breakPairingID)
				require.NotNil(t, cap.created)
			},
		},
		{
			name:     "通常系_旧ペアあり",
			menteeID: validMenteeID,
			mentorID: validMentorID,
			wantErr:  false,
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository, cap *pairingCreatorCaptures) {
				old := newTestPairing(oldPairingID, validMenteeID, 99)
				pairing.EXPECT().GetByMenteeID(validMenteeID).Return(old, nil)
				pairing.EXPECT().CreatePairing(validMenteeID, validMentorID).DoAndReturn(
					func(menteeID, mentorID uint) (*entity.Pairing, error) {
						cap.createMenteeID = menteeID
						cap.createMentorID = mentorID
						created := newTestPairing(202, menteeID, mentorID)
						cap.created = created
						return created, nil
					},
				)
				pairing.EXPECT().BreakPairing(oldPairingID).DoAndReturn(func(id uint) error {
					cap.breakPairingID = id
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *pairingCreatorCaptures) {
				t.Helper()
				assert.Equal(t, validMenteeID, cap.createMenteeID)
				assert.Equal(t, validMentorID, cap.createMentorID)
				assert.Equal(t, oldPairingID, cap.breakPairingID)
			},
		},
		{
			name:        "異常系_旧ペア取得失敗",
			menteeID:    validMenteeID,
			mentorID:    validMentorID,
			wantErr:     true,
			errContains: "ペアリング（旧）の取得に失敗しました",
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository, _ *pairingCreatorCaptures) {
				pairing.EXPECT().GetByMenteeID(validMenteeID).Return(nil, errGetPairingFailed)
			},
			assertFn: func(t *testing.T, cap *pairingCreatorCaptures) {
				t.Helper()
				assert.Equal(t, uint(0), cap.createMenteeID)
				assert.Equal(t, uint(0), cap.breakPairingID)
			},
		},
		{
			name:        "異常系_作成失敗",
			menteeID:    validMenteeID,
			mentorID:    validMentorID,
			wantErr:     true,
			errContains: "ペアリングの作成に失敗しました",
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository, cap *pairingCreatorCaptures) {
				pairing.EXPECT().GetByMenteeID(validMenteeID).Return(nil, nil)
				pairing.EXPECT().CreatePairing(validMenteeID, validMentorID).DoAndReturn(
					func(menteeID, mentorID uint) (*entity.Pairing, error) {
						cap.createMenteeID = menteeID
						cap.createMentorID = mentorID
						return nil, errCreatePairingFailed
					},
				)
			},
			assertFn: func(t *testing.T, cap *pairingCreatorCaptures) {
				t.Helper()
				assert.Equal(t, validMenteeID, cap.createMenteeID)
				assert.Equal(t, uint(0), cap.breakPairingID)
			},
		},
		{
			name:        "異常系_旧ペア解消失敗",
			menteeID:    validMenteeID,
			mentorID:    validMentorID,
			wantErr:     true,
			errContains: "ペアリングの解消に失敗しました",
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository, cap *pairingCreatorCaptures) {
				old := newTestPairing(oldPairingID, validMenteeID, 99)
				pairing.EXPECT().GetByMenteeID(validMenteeID).Return(old, nil)
				pairing.EXPECT().CreatePairing(validMenteeID, validMentorID).Return(newTestPairing(203, validMenteeID, validMentorID), nil)
				pairing.EXPECT().BreakPairing(oldPairingID).DoAndReturn(func(id uint) error {
					cap.breakPairingID = id
					return errBreakPairingFailed
				})
			},
			assertFn: func(t *testing.T, cap *pairingCreatorCaptures) {
				t.Helper()
				assert.Equal(t, oldPairingID, cap.breakPairingID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, pairingRepo := newPairingCreatorUsecase(t)
			cap := &pairingCreatorCaptures{}
			if tt.setup != nil {
				tt.setup(t, pairingRepo, cap)
			}

			gotErr := uc.PairingCreator(tt.menteeID, tt.mentorID)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "PairingCreator")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}
