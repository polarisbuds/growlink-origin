package planning_usecase

import (
	"testing"

	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"
	"growlink/internal/domain/valueobject"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestReferableUserIDs_AdminUnrestricted(t *testing.T) {
	uc := &PlanningAccessUsecase{}
	ids, err := uc.ReferableUserIDs(1, valueobject.SystemRoleAdmin)
	require.NoError(t, err)
	require.Nil(t, ids)
}

func TestReferableUserIDs_MentorMentee(t *testing.T) {
	ctrl := gomock.NewController(t)
	pairing := repomock.NewMockPairingRepository(ctrl)
	pairing.EXPECT().GetAll().Return([]entity.Pairing{
		{Model: gorm.Model{ID: 1}, MentorID: 9, MenteeID: 42, IsTemporary: true},
	}, nil)

	uc := &PlanningAccessUsecase{pairingRepo: pairing}
	ids, err := uc.ReferableUserIDs(9, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.True(t, ids[42])
}

func TestHasPairManagementAccess_NoScope(t *testing.T) {
	ctrl := gomock.NewController(t)
	teamMember := repomock.NewMockTeamMemberRepository(ctrl)
	teamMember.EXPECT().GetByMemberID(uint(7)).Return(nil, nil)

	uc := &PlanningAccessUsecase{teamMemberRepo: teamMember}

	ok, err := uc.HasPairManagementAccess(7, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.False(t, ok)
}
