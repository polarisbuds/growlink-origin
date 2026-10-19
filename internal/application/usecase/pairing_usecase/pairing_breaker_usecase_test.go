package pairing_usecase_test

import (
	"testing"

	"growlink/internal/application/usecase/pairing_usecase"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type pairingBreakerCaptures struct {
	breakPairingID uint
}

type pairingBreakerSetupFn func(t *testing.T, pairing *repomock.MockPairingRepository, cap *pairingBreakerCaptures)
type pairingBreakerAssertFn func(t *testing.T, cap *pairingBreakerCaptures)

func newPairingBreakerUsecase(t *testing.T) (
	*pairing_usecase.PairingBreakerUsecase,
	*repomock.MockPairingRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	pairingRepo := repomock.NewMockPairingRepository(ctrl)
	userRepo := repomock.NewMockUserRepository(ctrl)
	return pairing_usecase.NewPairingBreakerUsecase(pairingRepo, userRepo), pairingRepo
}

func TestPairingBreakerUsecase_PairingCreator(t *testing.T) {
	tests := []struct {
		name        string
		menteeID    uint
		wantErr     bool
		errContains string
		setup       pairingBreakerSetupFn
		assertFn    pairingBreakerAssertFn
	}{
		{
			name:     "通常系_解消",
			menteeID: validMenteeID,
			wantErr:  false,
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository, cap *pairingBreakerCaptures) {
				found := newTestPairing(oldPairingID, validMenteeID, validMentorID)
				pairing.EXPECT().GetByMenteeID(validMenteeID).Return(found, nil)
				pairing.EXPECT().BreakPairing(oldPairingID).DoAndReturn(func(id uint) error {
					cap.breakPairingID = id
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *pairingBreakerCaptures) {
				t.Helper()
				assert.Equal(t, oldPairingID, cap.breakPairingID)
			},
		},
		{
			name:        "異常系_取得失敗",
			menteeID:    validMenteeID,
			wantErr:     true,
			errContains: "削除対象のペアリング情報の取得に失敗しました",
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository, _ *pairingBreakerCaptures) {
				pairing.EXPECT().GetByMenteeID(validMenteeID).Return(nil, errGetPairingFailed)
			},
			assertFn: func(t *testing.T, cap *pairingBreakerCaptures) {
				t.Helper()
				assert.Equal(t, uint(0), cap.breakPairingID)
			},
		},
		{
			name:        "異常系_削除失敗",
			menteeID:    validMenteeID,
			wantErr:     true,
			errContains: "ペアリングの削除に失敗しました",
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository, cap *pairingBreakerCaptures) {
				found := newTestPairing(oldPairingID, validMenteeID, validMentorID)
				pairing.EXPECT().GetByMenteeID(validMenteeID).Return(found, nil)
				pairing.EXPECT().BreakPairing(oldPairingID).DoAndReturn(func(id uint) error {
					cap.breakPairingID = id
					return errBreakPairingFailed
				})
			},
			assertFn: func(t *testing.T, cap *pairingBreakerCaptures) {
				t.Helper()
				assert.Equal(t, oldPairingID, cap.breakPairingID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, pairingRepo := newPairingBreakerUsecase(t)
			cap := &pairingBreakerCaptures{}
			if tt.setup != nil {
				tt.setup(t, pairingRepo, cap)
			}

			got := uc.BreakPair(tt.menteeID)
			if tt.wantErr {
				assert.False(t, got)
			} else {
				assert.True(t, got)
			}
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}
