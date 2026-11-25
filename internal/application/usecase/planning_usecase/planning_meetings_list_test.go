package planning_usecase

import (
	"testing"
	"time"

	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func meetingOut(id uint, date time.Time, mentor, mentee string) outputmodel.MeetingOutput {
	return outputmodel.MeetingOutput{
		ID:            id,
		UUID:          uuid.New().String(),
		ImplementedAt: date,
		Status:        valueobject.MeetingStatusPending,
		Pairing: outputmodel.PairingOutput{
			Mentor: outputmodel.UserOutput{Name: mentor},
			Mentee: outputmodel.UserOutput{Name: mentee},
		},
	}
}

func TestSortPlanningMeetings_DateDescAndZeroLast(t *testing.T) {
	d1 := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	d2 := time.Date(2025, 1, 5, 12, 0, 0, 0, time.UTC)
	list := []outputmodel.MeetingOutput{
		meetingOut(1, time.Time{}, "a", "a"),
		meetingOut(2, d2, "a", "a"),
		meetingOut(3, d1, "a", "a"),
	}
	sortPlanningMeetings(list, true)
	require.Equal(t, uint(3), list[0].ID)
	require.Equal(t, uint(2), list[1].ID)
	require.Equal(t, uint(1), list[2].ID)
}

func TestSortPlanningMeetings_SameDateByMentorMenteeNameAsc(t *testing.T) {
	d := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	list := []outputmodel.MeetingOutput{
		meetingOut(1, d, "alpha", "m1"),
		meetingOut(2, d, "beta", "m1"),
		meetingOut(3, d, "beta", "m2"),
	}
	sortPlanningMeetings(list, true)
	require.Equal(t, "alpha", list[0].Pairing.Mentor.Name)
	require.Equal(t, "beta", list[1].Pairing.Mentor.Name)
	require.Equal(t, "m1", list[1].Pairing.Mentee.Name)
	require.Equal(t, uint(2), list[1].ID)
	require.Equal(t, uint(3), list[2].ID)
}

func TestSortPlanningMeetings_TieBreakByIDDESC(t *testing.T) {
	d := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	list := []outputmodel.MeetingOutput{
		meetingOut(10, d, "同", "同"),
		meetingOut(20, d, "同", "同"),
	}
	sortPlanningMeetings(list, true)
	require.Equal(t, uint(20), list[0].ID)
	require.Equal(t, uint(10), list[1].ID)
}

func TestSlicePlanningMeetingsPage_OutOfRange(t *testing.T) {
	list := make([]outputmodel.MeetingOutput, 25)
	for i := range list {
		list[i] = meetingOut(uint(i+1), time.Now(), "a", "b")
	}
	page := slicePlanningMeetingsPage(list, 2, PlanningMeetingsPageSize)
	require.Len(t, page, 5)
	empty := slicePlanningMeetingsPage(list, 10, PlanningMeetingsPageSize)
	require.Empty(t, empty)
}