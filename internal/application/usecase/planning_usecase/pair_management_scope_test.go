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

func TestPairManagementScopeUserIDs_NoOrgRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	teamMember := repomock.NewMockTeamMemberRepository(ctrl)
	teamMember.EXPECT().GetByMemberID(uint(7)).Return(nil, nil)

	uc := &PlanningAccessUsecase{teamMemberRepo: teamMember}

	ids, err := uc.PairManagementScopeUserIDs(7, valueobject.AuthorityGeneral, valueobject.SystemRoleNone)
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

	ids, err := uc.PairManagementScopeUserIDs(3, valueobject.AuthorityGeneral, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.True(t, ids[42])
	require.True(t, ids[3])
}

func TestPairManagementScopeUserIDs_GeneralMentorMentees(t *testing.T) {
	ctrl := gomock.NewController(t)
	teamMember := repomock.NewMockTeamMemberRepository(ctrl)
	teamMember.EXPECT().GetByMemberID(uint(9)).Return(nil, nil)
	pairing := repomock.NewMockPairingRepository(ctrl)
	pairing.EXPECT().GetByMentorID(uint(9)).Return([]entity.Pairing{
		{Model: gorm.Model{ID: 1}, MentorID: 9, MenteeID: 42},
	}, nil)

	uc := &PlanningAccessUsecase{teamMemberRepo: teamMember, pairingRepo: pairing}

	ids, err := uc.PairManagementScopeUserIDs(9, valueobject.AuthorityGeneral, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.True(t, ids[42])
}

func TestPairManagementScopeUserIDs_BossWithoutDeptRole_NoSubordinates(t *testing.T) {
	ctrl := gomock.NewController(t)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)
	deptMember.EXPECT().GetByUserID(uint(5)).Return([]entity.DepartmentMember{
		{DepartmentID: 1, UserID: 5, IsHead: false, IsManager: false},
	}, nil)
	hierarchy := repomock.NewMockDepartmentHierarchyRepository(ctrl)
	teamMember := repomock.NewMockTeamMemberRepository(ctrl)
	teamMember.EXPECT().GetByMemberID(uint(5)).Return(nil, nil)

	uc := &PlanningAccessUsecase{
		deptMemberRepo: deptMember,
		hierarchyRepo:  hierarchy,
		teamMemberRepo: teamMember,
	}

	ids, err := uc.PairManagementScopeUserIDs(5, valueobject.AuthorityBoss, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.Empty(t, ids)
}

func TestPairManagementScopeUserIDs_BossSkipsMentorMentees(t *testing.T) {
	ctrl := gomock.NewController(t)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)
	deptMember.EXPECT().GetByUserID(uint(9)).Return(nil, nil)
	hierarchy := repomock.NewMockDepartmentHierarchyRepository(ctrl)
	teamMember := repomock.NewMockTeamMemberRepository(ctrl)
	teamMember.EXPECT().GetByMemberID(uint(9)).Return(nil, nil)

	uc := &PlanningAccessUsecase{
		deptMemberRepo: deptMember,
		hierarchyRepo:  hierarchy,
		teamMemberRepo: teamMember,
		pairingRepo:    repomock.NewMockPairingRepository(ctrl),
	}

	ids, err := uc.PairManagementScopeUserIDs(9, valueobject.AuthorityBoss, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.Empty(t, ids)
}

func TestHasPairManagementAccess_UsesPairScopeNotMentor(t *testing.T) {
	ctrl := gomock.NewController(t)
	teamMember := repomock.NewMockTeamMemberRepository(ctrl)
	teamMember.EXPECT().GetByMemberID(uint(9)).Return(nil, nil)
	pairing := repomock.NewMockPairingRepository(ctrl)
	pairing.EXPECT().GetByMentorID(uint(9)).Return(nil, nil)

	uc := &PlanningAccessUsecase{teamMemberRepo: teamMember, pairingRepo: pairing}

	ok, err := uc.HasPairManagementAccess(9, valueobject.AuthorityGeneral, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.False(t, ok)
}

func TestHasPairManagementAccess_UsesPairScopeMentorWithMentee(t *testing.T) {
	ctrl := gomock.NewController(t)
	teamMember := repomock.NewMockTeamMemberRepository(ctrl)
	teamMember.EXPECT().GetByMemberID(uint(9)).Return(nil, nil)
	pairing := repomock.NewMockPairingRepository(ctrl)
	pairing.EXPECT().GetByMentorID(uint(9)).Return([]entity.Pairing{
		{Model: gorm.Model{ID: 1}, MentorID: 9, MenteeID: 42},
	}, nil)

	uc := &PlanningAccessUsecase{teamMemberRepo: teamMember, pairingRepo: pairing}

	ok, err := uc.HasPairManagementAccess(9, valueobject.AuthorityGeneral, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestHasPairManagementAccess_Authority_Owner(t *testing.T) {
	uc := &PlanningAccessUsecase{}

	ok, err := uc.HasPairManagementAccess(9, valueobject.AuthorityOwner, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.True(t, ok)
}
