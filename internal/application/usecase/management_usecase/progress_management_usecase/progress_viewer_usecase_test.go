package progress_mng_usecase

import (
	"errors"
	"testing"

	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"
	"growlink/internal/domain/valueobject"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type stubReferable struct {
	ids map[uint]bool
	err error
}

func (s *stubReferable) ReferableUserIDs(uint, valueobject.SystemRole) (map[uint]bool, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.ids, nil
}

func meetingForMentee(menteeID uint, status valueobject.MeetingStatus) entity.Meeting {
	return entity.Meeting{
		Status:  status,
		Pairing: &entity.Pairing{MenteeID: menteeID},
	}
}

func TestProgressViewer_Summary_AdminSeesAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	meeting := repomock.NewMockMeetingRepository(ctrl)
	meeting.EXPECT().GetAll().Return([]entity.Meeting{
		meetingForMentee(1, valueobject.MeetingStatusPending),
		meetingForMentee(2, valueobject.MeetingStatusDone),
	}, nil)
	uc := NewProgressViewer(meeting, &stubReferable{ids: map[uint]bool{99: true}})

	got, err := uc.Summary(ProgressFilter{ViewerID: 9, SystemRole: valueobject.SystemRoleAdmin})
	require.NoError(t, err)
	require.Equal(t, 2, got.Total)
}

func TestProgressViewer_Summary_AnalystSeesAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	meeting := repomock.NewMockMeetingRepository(ctrl)
	meeting.EXPECT().GetAll().Return([]entity.Meeting{
		meetingForMentee(1, valueobject.MeetingStatusPending),
	}, nil)
	uc := NewProgressViewer(meeting, &stubReferable{})

	got, err := uc.Summary(ProgressFilter{ViewerID: 9, SystemRole: valueobject.SystemRoleAnalyst})
	require.NoError(t, err)
	require.Equal(t, 1, got.Total)
}

func TestProgressViewer_Summary_ReferableScopeOnly(t *testing.T) {
	ctrl := gomock.NewController(t)
	meeting := repomock.NewMockMeetingRepository(ctrl)
	meeting.EXPECT().GetAll().Return([]entity.Meeting{
		meetingForMentee(10, valueobject.MeetingStatusPending),
		meetingForMentee(99, valueobject.MeetingStatusDone),
	}, nil)
	uc := NewProgressViewer(meeting, &stubReferable{ids: map[uint]bool{10: true}})

	got, err := uc.Summary(ProgressFilter{ViewerID: 9, SystemRole: valueobject.SystemRoleNone})
	require.NoError(t, err)
	require.Equal(t, 1, got.Total)
	require.Equal(t, 1, got.Pending)
}

func TestProgressViewer_Summary_NoReferableUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	meeting := repomock.NewMockMeetingRepository(ctrl)
	meeting.EXPECT().GetAll().Return([]entity.Meeting{
		meetingForMentee(1, valueobject.MeetingStatusPending),
	}, nil)
	uc := NewProgressViewer(meeting, &stubReferable{ids: map[uint]bool{}})

	got, err := uc.Summary(ProgressFilter{ViewerID: 9, SystemRole: valueobject.SystemRoleNone})
	require.NoError(t, err)
	require.Equal(t, 0, got.Total)
}

func TestProgressViewer_Summary_ReferableError(t *testing.T) {
	ctrl := gomock.NewController(t)
	meeting := repomock.NewMockMeetingRepository(ctrl)
	meeting.EXPECT().GetAll().Return([]entity.Meeting{
		meetingForMentee(1, valueobject.MeetingStatusPending),
	}, nil)
	uc := NewProgressViewer(meeting, &stubReferable{err: errors.New("scope failed")})

	_, err := uc.Summary(ProgressFilter{ViewerID: 9, SystemRole: valueobject.SystemRoleNone})
	require.Error(t, err)
}
