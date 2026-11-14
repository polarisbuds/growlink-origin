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
		meetingOut(1, d, "山田", "あ"),
		meetingOut(2, d, "佐藤", "あ"),
		meetingOut(3, d, "佐藤", "い"),
	}
	sortPlanningMeetings(list, true)
	require.Equal(t, "山田", list[0].Pairing.Mentor.Name)
	require.Equal(t, "佐藤", list[1].Pairing.Mentor.Name)
	require.Equal(t, "い", list[1].Pairing.Mentee.Name)
	require.Equal(t, uint(3), list[1].ID)
	require.Equal(t, uint(2), list[2].ID)
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

func TestFilterPlanningMeetings_DateRangeEmpty(t *testing.T) {
	f := PlanningMeetingsListFilter{DateRangeEmpty: true}
	out := filterPlanningMeetings([]outputmodel.MeetingOutput{meetingOut(1, time.Now(), "a", "b")}, f, nil)
	require.Empty(t, out)
}

func TestFilterPlanningMeetings_ExcludeZeroDateWhenFromSet(t *testing.T) {
	from, err := JSTDayStart("2025-01-01")
	require.NoError(t, err)
	f := PlanningMeetingsListFilter{From: &from}
	list := []outputmodel.MeetingOutput{
		meetingOut(1, time.Time{}, "a", "b"),
		meetingOut(2, time.Date(2025, 1, 2, 0, 0, 0, 0, jst), "a", "b"),
	}
	out := filterPlanningMeetings(list, f, nil)
	require.Len(t, out, 1)
	require.Equal(t, uint(2), out[0].ID)
}

func TestFilterPlanningMeetings_Status(t *testing.T) {
	st := valueobject.MeetingStatusDone
	f := PlanningMeetingsListFilter{Status: &st}
	m := meetingOut(1, time.Now(), "a", "b")
	m.Status = valueobject.MeetingStatusDone
	m2 := meetingOut(2, time.Now(), "a", "b")
	m2.Status = valueobject.MeetingStatusDraft
	out := filterPlanningMeetings([]outputmodel.MeetingOutput{m, m2}, f, nil)
	require.Len(t, out, 1)
	require.Equal(t, uint(1), out[0].ID)
}

func TestFilterPlanningMeetings_DepartmentOrMentorMentee(t *testing.T) {
	f := PlanningMeetingsListFilter{DepartmentUUIDs: []uuid.UUID{uuid.New()}}
	m1 := meetingOut(1, time.Now(), "a", "b")
	m1.Pairing.Mentor.UUID = "mentor-in"
	m2 := meetingOut(2, time.Now(), "c", "d")
	m2.Pairing.Mentee.UUID = "mentee-in"
	deptUsers := map[string]struct{}{"mentor-in": {}, "mentee-in": {}}
	out := filterPlanningMeetings([]outputmodel.MeetingOutput{m1, m2}, f, deptUsers)
	require.Len(t, out, 2)
}

func TestSlicePlanningMeetingsPage_OutOfRange(t *testing.T) {
	list := make([]outputmodel.MeetingOutput, 25)
	for i := range list {
		list[i] = meetingOut(uint(i+1), time.Now(), "a", "b")
	}
	page := slicePlanningMeetingsPage(list, 3, PlanningMeetingsPageSize)
	require.Len(t, page, 5)
	empty := slicePlanningMeetingsPage(list, 10, PlanningMeetingsPageSize)
	require.Empty(t, empty)
}