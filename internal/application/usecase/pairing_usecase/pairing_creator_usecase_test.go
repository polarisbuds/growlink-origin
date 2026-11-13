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
	createMenteeID        uint
	createMentorID        uint
	createCreatedByUserID uint
	created               *entity.Pairing
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
			name:     "通常系_作成",
			menteeID: validMenteeID,
			mentorID: validMentorID,
			wantErr:  false,
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository, cap *pairingCreatorCaptures) {
				pairing.EXPECT().CreatePairing(validMenteeID, validMentorID, uint(0)).DoAndReturn(
					func(menteeID, mentorID, createdByUserID uint) (*entity.Pairing, error) {
						cap.createMenteeID = menteeID
						cap.createMentorID = mentorID
						cap.createCreatedByUserID = createdByUserID
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
				assert.Equal(t, uint(0), cap.createCreatedByUserID)
				require.NotNil(t, cap.created)
			},
		},
		{
			name:     "通常系_既存ペアがあっても追加のみ",
			menteeID: validMenteeID,
			mentorID: validMentorID,
			wantErr:  false,
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository, cap *pairingCreatorCaptures) {
				pairing.EXPECT().CreatePairing(validMenteeID, validMentorID, uint(0)).DoAndReturn(
					func(menteeID, mentorID, createdByUserID uint) (*entity.Pairing, error) {
						cap.createMenteeID = menteeID
						cap.createMentorID = mentorID
						cap.createCreatedByUserID = createdByUserID
						created := newTestPairing(202, menteeID, mentorID)
						cap.created = created
						return created, nil
					},
				)
			},
			assertFn: func(t *testing.T, cap *pairingCreatorCaptures) {
				t.Helper()
				assert.Equal(t, validMenteeID, cap.createMenteeID)
				assert.Equal(t, validMentorID, cap.createMentorID)
			},
		},
		{
			name:        "異常系_作成失敗",
			menteeID:    validMenteeID,
			mentorID:    validMentorID,
			wantErr:     true,
			errContains: "ペアリングの作成に失敗しました",
			setup: func(_ *testing.T, pairing *repomock.MockPairingRepository, cap *pairingCreatorCaptures) {
				pairing.EXPECT().CreatePairing(validMenteeID, validMentorID, uint(0)).DoAndReturn(
					func(menteeID, mentorID, createdByUserID uint) (*entity.Pairing, error) {
						cap.createMenteeID = menteeID
						cap.createMentorID = mentorID
						cap.createCreatedByUserID = createdByUserID
						return nil, errCreatePairingFailed
					},
				)
			},
			assertFn: func(t *testing.T, cap *pairingCreatorCaptures) {
				t.Helper()
				assert.Equal(t, validMenteeID, cap.createMenteeID)
				assert.Equal(t, uint(0), cap.createCreatedByUserID)
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
