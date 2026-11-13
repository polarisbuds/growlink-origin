package progress_mng_usecase

import (
	"testing"
	"time"

	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func newDeptProgressUC(t *testing.T, referable ReferableUsers, meetings []entity.Meeting, pairings []entity.Pairing) *ProgressViewer {
	t.Helper()
	ctrl := gomock.NewController(t)
	meeting := repomock.NewMockMeetingRepository(ctrl)
	meeting.EXPECT().GetAll().Return(meetings, nil).AnyTimes()
	pairing := repomock.NewMockPairingRepository(ctrl)
	pairing.EXPECT().GetAll().Return(pairings, nil).AnyTimes()
	dept := repomock.NewMockDepartmentRepository(ctrl)
	dept.EXPECT().GetAll().Return([]entity.Department{
		{Model: gorm.Model{ID: 1}, UUID: uuid.New(), Name: "本部"},
		{Model: gorm.Model{ID: 2}, UUID: uuid.New(), Name: "第1事業部"},
	}, nil).AnyTimes()
	hier := repomock.NewMockDepartmentHierarchyRepository(ctrl)
	hier.EXPECT().GetAll().Return([]entity.DepartmentHierarchy{
		{ParentDepartmentID: 1, ChildDepartmentID: 2},
	}, nil).AnyTimes()
	deptMem := repomock.NewMockDepartmentMemberRepository(ctrl)
	deptMem.EXPECT().GetAll().Return([]entity.DepartmentMember{
		{DepartmentID: 1, UserID: 10},
		{DepartmentID: 2, UserID: 20},
		{DepartmentID: 2, UserID: 21},
	}, nil).AnyTimes()
	team := repomock.NewMockTeamRepository(ctrl)
	team.EXPECT().GetAll().Return([]entity.Team{
		{Model: gorm.Model{ID: 100}, DepartmentID: 2, Name: "Aチーム"},
	}, nil).AnyTimes()
	teamMem := repomock.NewMockTeamMemberRepository(ctrl)
	teamMem.EXPECT().GetAll().Return([]entity.TeamMember{
		{TeamID: 100, MemberID: 30},
	}, nil).AnyTimes()
	return NewProgressViewer(meeting, referable, pairing, dept, hier, deptMem, team, teamMem)
}

func meetingMentee(menteeID uint, status valueobject.MeetingStatus, at time.Time) entity.Meeting {
	return entity.Meeting{
		Status:        status,
		ImplementedAt: &at,
		Pairing:       &entity.Pairing{MenteeID: menteeID},
	}
}

func findDeptRow(rows []outputmodel.ProgressDepartmentRowOutput, name string) *outputmodel.ProgressDepartmentRowOutput {
	for i := range rows {
		if rows[i].Name == name {
			return &rows[i]
		}
	}
	return nil
}

func TestDepartmentTable_IncludesSubtreeAndTeam(t *testing.T) {
	now := time.Date(2026, 11, 15, 0, 0, 0, 0, time.Local)
	meetings := []entity.Meeting{
		meetingMentee(20, valueobject.MeetingStatusPending, now),
		meetingMentee(30, valueobject.MeetingStatusDone, now),
	}
	pairings := []entity.Pairing{
		{MenteeID: 20},
		{MenteeID: 30},
	}
	uc := newDeptProgressUC(t, &stubReferable{}, meetings, pairings)
	got, err := uc.DepartmentTable(ProgressFilter{SystemRole: valueobject.SystemRoleAdmin})
	require.NoError(t, err)
	require.NotEmpty(t, got.Rows)
	require.True(t, got.Rows[0].IsTotal)
	require.GreaterOrEqual(t, got.Rows[0].UserCount, 4)
	child := findDeptRow(got.Rows, "第1事業部")
	require.NotNil(t, child)
	require.Equal(t, 3, child.UserCount)
	require.Equal(t, 1, child.NotDoneCount)
	require.Equal(t, 1, child.DoneCount)
	head := findDeptRow(got.Rows, "本部")
	require.NotNil(t, head)
	require.Empty(t, head.Teams)
	require.Len(t, child.Teams, 1)
	require.Equal(t, "Aチーム", child.Teams[0].Name)
	require.Equal(t, 1, child.Teams[0].UserCount)
	require.Equal(t, 1, child.Teams[0].DoneCount)
}

func TestDepartmentTable_TeamRowRespectsVisibility(t *testing.T) {
	now := time.Date(2026, 11, 15, 0, 0, 0, 0, time.Local)
	meetings := []entity.Meeting{meetingMentee(30, valueobject.MeetingStatusDone, now)}
	pairings := []entity.Pairing{{MenteeID: 30}}
	uc := newDeptProgressUC(t, &stubReferable{ids: map[uint]bool{20: true}}, meetings, pairings)
	got, err := uc.DepartmentTable(ProgressFilter{ViewerID: 5, SystemRole: valueobject.SystemRoleNone})
	require.NoError(t, err)
	child := findDeptRow(got.Rows, "第1事業部")
	require.NotNil(t, child)
	require.Empty(t, child.Teams)
}

func TestDepartmentTable_VisibleScopeMentorOnly(t *testing.T) {
	now := time.Date(2026, 11, 15, 0, 0, 0, 0, time.Local)
	meetings := []entity.Meeting{meetingMentee(20, valueobject.MeetingStatusSubmitted, now)}
	pairings := []entity.Pairing{{MenteeID: 20}}
	uc := newDeptProgressUC(t, &stubReferable{ids: map[uint]bool{20: true}}, meetings, pairings)
	got, err := uc.DepartmentTable(ProgressFilter{ViewerID: 5, SystemRole: valueobject.SystemRoleNone})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(got.Rows), 2)
	require.Equal(t, 1, got.Rows[0].UserCount)
	require.Equal(t, 1, got.Rows[0].SubmittedCount)
	child := findDeptRow(got.Rows, "第1事業部")
	require.NotNil(t, child)
	require.Equal(t, 1, child.SubmittedCount)
}

func TestVisibleUserIDs_EmptyReferableAddsViewer(t *testing.T) {
	uc := &ProgressViewer{referable: &stubReferable{ids: map[uint]bool{}}}
	got, err := uc.visibleUserIDs(ProgressFilter{ViewerID: 42, SystemRole: valueobject.SystemRoleNone})
	require.NoError(t, err)
	require.True(t, got[42])
}
