package planning_usecase

import (
	"testing"

	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"
	"growlink/internal/domain/valueobject"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPairManagementScopeUserIDs_NoOrgRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	teamMember := repomock.NewMockTeamMemberRepository(ctrl)
	teamMember.EXPECT().GetByMemberID(uint(7)).Return(nil, nil)

	uc := &PlanningAccessUsecase{teamMemberRepo: teamMember}

	ids, err := uc.PairManagementScopeUserIDs(7, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.Empty(t, ids)
}

func TestPairManagementScopeUserIDs_TeamLeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	teamMember := repomock.NewMockTeamMemberRepository(ctrl)
	teamMember.EXPECT().GetByMemberID(uint(3)).Return([]entity.TeamMember{
		{TeamID: 10, MemberID: 3, IsReader: true},
	}, nil)
	teamMember.EXPECT().GetByTeamID(uint(10)).Return([]entity.TeamMember{
		{TeamID: 10, MemberID: 3, IsReader: true},
		{TeamID: 10, MemberID: 42, IsReader: false},
	}, nil)

	uc := &PlanningAccessUsecase{teamMemberRepo: teamMember}

	ids, err := uc.PairManagementScopeUserIDs(3, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.True(t, ids[42])
	require.True(t, ids[3])
}

func TestHasPairManagementAccess_UsesPairScopeNotMentor(t *testing.T) {
	ctrl := gomock.NewController(t)
	teamMember := repomock.NewMockTeamMemberRepository(ctrl)
	teamMember.EXPECT().GetByMemberID(uint(9)).Return(nil, nil)

	uc := &PlanningAccessUsecase{teamMemberRepo: teamMember}

	ok, err := uc.HasPairManagementAccess(9, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.False(t, ok)
}
