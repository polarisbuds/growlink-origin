package progress_mng_usecase

import (
	"errors"
	"testing"

	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"
	"growlink/internal/domain/valueobject"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type progressViewerMocks struct {
	meeting      *repomock.MockMeetingRepository
	user         *repomock.MockUserRepository
	teamMember   *repomock.MockTeamMemberRepository
	userDept     *stubUserDeptRepo
	deptRef      *repomock.MockUserReferenceDepartmentRepository
	teamRef      *repomock.MockUserReferenceTeamRepository
	relationship *repomock.MockUserRelationshipRepository
}

func newProgressViewer(t *testing.T) (*ProgressViewer, *progressViewerMocks) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mocks := &progressViewerMocks{
		meeting:      repomock.NewMockMeetingRepository(ctrl),
		user:         repomock.NewMockUserRepository(ctrl),
		teamMember:   repomock.NewMockTeamMemberRepository(ctrl),
		userDept:     &stubUserDeptRepo{byDept: map[uint][]uint{}},
		deptRef:      repomock.NewMockUserReferenceDepartmentRepository(ctrl),
		teamRef:      repomock.NewMockUserReferenceTeamRepository(ctrl),
		relationship: repomock.NewMockUserRelationshipRepository(ctrl),
	}
	uc := NewProgressViewer(
		mocks.meeting,
		mocks.user,
		mocks.teamMember,
		mocks.userDept,
		mocks.deptRef,
		mocks.teamRef,
		mocks.relationship,
	)
	return uc, mocks
}

func meetingForMentee(menteeID uint, status valueobject.MeetingStatus) entity.Meeting {
	return entity.Meeting{
		Status:  status,
		Pairing: &entity.Pairing{MenteeID: menteeID},
	}
}

func TestProgressViewer_Summary_AdminSeesAll(t *testing.T) {
	uc, mocks := newProgressViewer(t)
	mocks.meeting.EXPECT().GetAll().Return([]entity.Meeting{
		meetingForMentee(1, valueobject.MeetingStatusPending),
		meetingForMentee(2, valueobject.MeetingStatusDone),
		meetingForMentee(3, valueobject.MeetingStatusDraft),
	}, nil)

	got, err := uc.Summary(ProgressFilter{ViewerID: 9, SystemRole: valueobject.SystemRoleAdmin})
	require.NoError(t, err)
	require.Equal(t, 3, got.Total)
	require.Equal(t, 1, got.Pending)
	require.Equal(t, 1, got.Draft)
	require.Equal(t, 1, got.Done)
}

func TestProgressViewer_Summary_AnalystSeesAll(t *testing.T) {
	uc, mocks := newProgressViewer(t)
	mocks.meeting.EXPECT().GetAll().Return([]entity.Meeting{
		meetingForMentee(1, valueobject.MeetingStatusPending),
		meetingForMentee(2, valueobject.MeetingStatusPending),
	}, nil)

	got, err := uc.Summary(ProgressFilter{ViewerID: 9, SystemRole: valueobject.SystemRoleAnalyst})
	require.NoError(t, err)
	require.Equal(t, 2, got.Total)
}

func TestProgressViewer_Summary_NoneSeesNothing(t *testing.T) {
	uc, mocks := newProgressViewer(t)
	mocks.user.EXPECT().GetByID(uint(9)).Return(&entity.User{
		OperationTarget: valueobject.OperationTargetNone,
	}, nil)
	mocks.meeting.EXPECT().GetAll().Return([]entity.Meeting{
		meetingForMentee(1, valueobject.MeetingStatusPending),
		meetingForMentee(2, valueobject.MeetingStatusDone),
	}, nil)

	got, err := uc.Summary(ProgressFilter{ViewerID: 9, SystemRole: valueobject.SystemRoleNone})
	require.NoError(t, err)
	require.Equal(t, 0, got.Total)
}

func TestProgressViewer_Summary_DepartmentScope(t *testing.T) {
	uc, mocks := newProgressViewer(t)
	mocks.user.EXPECT().GetByID(uint(9)).Return(&entity.User{
		OperationTarget: valueobject.OperationTargetDepartment,
	}, nil)
	mocks.deptRef.EXPECT().GetByUserID(uint(9)).Return([]entity.UserReferenceDepartment{
		{ReferableDepartmentID: 10},
	}, nil)
	mocks.user.EXPECT().GetAll().Return([]entity.User{
		{Model: gorm.Model{ID: 1}, DepartmentID: 10},
		{Model: gorm.Model{ID: 2}, DepartmentID: 11},
		{Model: gorm.Model{ID: 4}, DepartmentID: 99},
	}, nil)
	mocks.userDept.byDept[10] = []uint{3}
	mocks.meeting.EXPECT().GetAll().Return([]entity.Meeting{
		meetingForMentee(1, valueobject.MeetingStatusPending),
		meetingForMentee(2, valueobject.MeetingStatusPending),
		meetingForMentee(3, valueobject.MeetingStatusDone),
		meetingForMentee(4, valueobject.MeetingStatusDraft),
	}, nil)

	got, err := uc.Summary(ProgressFilter{ViewerID: 9, SystemRole: valueobject.SystemRoleNone})
	require.NoError(t, err)
	require.Equal(t, 2, got.Total)
	require.Equal(t, 1, got.Pending)
	require.Equal(t, 1, got.Done)
	require.Equal(t, 0, got.Draft)
}

func TestProgressViewer_Summary_TeamScope(t *testing.T) {
	uc, mocks := newProgressViewer(t)
	mocks.user.EXPECT().GetByID(uint(9)).Return(&entity.User{
		OperationTarget: valueobject.OperationTargetTeam,
	}, nil)
	mocks.teamRef.EXPECT().GetByUserID(uint(9)).Return([]entity.UserReferenceTeam{
		{ReferableTeamID: 5},
	}, nil)
	mocks.teamMember.EXPECT().GetByTeamID(uint(5)).Return([]entity.TeamMember{
		{MemberID: 21},
		{MemberID: 22},
	}, nil)
	mocks.meeting.EXPECT().GetAll().Return([]entity.Meeting{
		meetingForMentee(21, valueobject.MeetingStatusPending),
		meetingForMentee(99, valueobject.MeetingStatusDone),
	}, nil)

	got, err := uc.Summary(ProgressFilter{ViewerID: 9, SystemRole: valueobject.SystemRoleNone})
	require.NoError(t, err)
	require.Equal(t, 1, got.Total)
	require.Equal(t, 1, got.Pending)
}

func TestProgressViewer_Summary_SubordinateScope(t *testing.T) {
	uc, mocks := newProgressViewer(t)
	mocks.user.EXPECT().GetByID(uint(9)).Return(&entity.User{
		OperationTarget: valueobject.OperationTargetSubordinate,
	}, nil)
	mocks.relationship.EXPECT().GetByBossID(uint(9)).Return([]entity.UserRelationship{
		{SourceUserID: 31},
	}, nil)
	mocks.meeting.EXPECT().GetAll().Return([]entity.Meeting{
		meetingForMentee(31, valueobject.MeetingStatusDraft),
		meetingForMentee(32, valueobject.MeetingStatusDone),
	}, nil)

	got, err := uc.Summary(ProgressFilter{ViewerID: 9, SystemRole: valueobject.SystemRoleNone})
	require.NoError(t, err)
	require.Equal(t, 1, got.Total)
	require.Equal(t, 1, got.Draft)
}

func TestProgressViewer_Summary_DepartmentLookupError(t *testing.T) {
	uc, mocks := newProgressViewer(t)
	mocks.user.EXPECT().GetByID(uint(9)).Return(&entity.User{
		OperationTarget: valueobject.OperationTargetDepartment,
	}, nil)
	mocks.deptRef.EXPECT().GetByUserID(uint(9)).Return(nil, errors.New("ref failed"))
	mocks.meeting.EXPECT().GetAll().Return([]entity.Meeting{
		meetingForMentee(1, valueobject.MeetingStatusPending),
	}, nil)

	_, err := uc.Summary(ProgressFilter{ViewerID: 9, SystemRole: valueobject.SystemRoleNone})
	require.Error(t, err)
}

type stubUserDeptRepo struct {
	byDept map[uint][]uint
}

func (s *stubUserDeptRepo) GetUserIDsByDepartmentID(departmentID uint) ([]uint, error) {
	return s.byDept[departmentID], nil
}

func (s *stubUserDeptRepo) SyncDepartmentMembers(uint, []uint) error { return nil }

func (s *stubUserDeptRepo) GetDepartmentIDsByUserID(uint) ([]uint, error) { return nil, nil }

func (s *stubUserDeptRepo) Create(uint, uint) error { return nil }

func (s *stubUserDeptRepo) DeleteByUserID(uint) error { return nil }

func (s *stubUserDeptRepo) DeleteByUserAndDepartment(uint, uint) error { return nil }
