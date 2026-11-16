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

func TestDepartmentScopeUserIDs_ManagerExcludesHeadOfRootDept(t *testing.T) {
	ctrl := gomock.NewController(t)
	const (
		deptID   = uint(1)
		manager  = uint(10)
		head     = uint(20)
		employee = uint(30)
	)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)
	deptMember.EXPECT().GetByUserID(manager).Return([]entity.DepartmentMember{
		{DepartmentID: deptID, UserID: manager, IsHead: false, IsManager: true},
	}, nil)
	hierarchy := repomock.NewMockDepartmentHierarchyRepository(ctrl)
	hierarchy.EXPECT().GetAll().Return(nil, nil)
	deptMember.EXPECT().GetByDepartmentID(deptID).Return([]entity.DepartmentMember{
		{DepartmentID: deptID, UserID: head, IsHead: true, IsManager: true},
		{DepartmentID: deptID, UserID: manager, IsHead: false, IsManager: true},
		{DepartmentID: deptID, UserID: employee, IsHead: false, IsManager: false},
	}, nil).Times(2)

	uc := &PlanningAccessUsecase{
		deptMemberRepo: deptMember,
		hierarchyRepo:  hierarchy,
	}

	ids, err := uc.departmentScopeUserIDs(manager)
	require.NoError(t, err)
	require.Contains(t, ids, employee)
	require.Contains(t, ids, manager)
	require.NotContains(t, ids, head)
}

func TestDepartmentScopeUserIDs_HeadIncludesManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	const (
		deptID  = uint(1)
		head    = uint(20)
		manager = uint(10)
	)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)
	deptMember.EXPECT().GetByUserID(head).Return([]entity.DepartmentMember{
		{DepartmentID: deptID, UserID: head, IsHead: true, IsManager: true},
	}, nil)
	hierarchy := repomock.NewMockDepartmentHierarchyRepository(ctrl)
	hierarchy.EXPECT().GetAll().Return(nil, nil)
	deptMember.EXPECT().GetByDepartmentID(deptID).Return([]entity.DepartmentMember{
		{DepartmentID: deptID, UserID: head, IsHead: true, IsManager: true},
		{DepartmentID: deptID, UserID: manager, IsHead: false, IsManager: true},
	}, nil)

	uc := &PlanningAccessUsecase{
		deptMemberRepo: deptMember,
		hierarchyRepo:  hierarchy,
	}

	ids, err := uc.departmentScopeUserIDs(head)
	require.NoError(t, err)
	require.Contains(t, ids, head)
	require.Contains(t, ids, manager)
}

func TestDepartmentScopeUserIDs_ManagerAndHeadInDifferentDepts(t *testing.T) {
	ctrl := gomock.NewController(t)
	const (
		managerDept = uint(1)
		headDept    = uint(2)
		user        = uint(100)
		headInMgr   = uint(20)
		subordinate = uint(40)
	)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)
	deptMember.EXPECT().GetByUserID(user).Return([]entity.DepartmentMember{
		{DepartmentID: managerDept, UserID: user, IsHead: false, IsManager: true},
		{DepartmentID: headDept, UserID: user, IsHead: true, IsManager: true},
	}, nil)
	hierarchy := repomock.NewMockDepartmentHierarchyRepository(ctrl)
	hierarchy.EXPECT().GetAll().Return(nil, nil)

	deptMember.EXPECT().GetByDepartmentID(managerDept).Return([]entity.DepartmentMember{
		{DepartmentID: managerDept, UserID: headInMgr, IsHead: true, IsManager: true},
		{DepartmentID: managerDept, UserID: user, IsHead: false, IsManager: true},
	}, nil).Times(2)
	deptMember.EXPECT().GetByDepartmentID(headDept).Return([]entity.DepartmentMember{
		{DepartmentID: headDept, UserID: user, IsHead: true, IsManager: true},
		{DepartmentID: headDept, UserID: subordinate, IsHead: false, IsManager: false},
	}, nil)

	uc := &PlanningAccessUsecase{
		deptMemberRepo: deptMember,
		hierarchyRepo:  hierarchy,
	}

	ids, err := uc.departmentScopeUserIDs(user)
	require.NoError(t, err)
	require.NotContains(t, ids, headInMgr)
	require.Contains(t, ids, subordinate)
	require.Contains(t, ids, user)
}

func TestHasPairManagementAccess_NoScope(t *testing.T) {
	ctrl := gomock.NewController(t)
	teamMember := repomock.NewMockTeamMemberRepository(ctrl)
	teamMember.EXPECT().GetByMemberID(uint(7)).Return(nil, nil)
	pairing := repomock.NewMockPairingRepository(ctrl)
	pairing.EXPECT().GetByMentorID(uint(7)).Return(nil, nil)

	uc := &PlanningAccessUsecase{teamMemberRepo: teamMember, pairingRepo: pairing}

	ok, err := uc.HasPairManagementAccess(7, valueobject.AuthorityGeneral, valueobject.SystemRoleNone)
	require.NoError(t, err)
	require.False(t, ok)
}

