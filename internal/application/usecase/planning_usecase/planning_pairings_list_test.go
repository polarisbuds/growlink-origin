package planning_usecase

import (
	"testing"

	outputmodel "growlink/internal/application/model/output_model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func pairingOut(id uint, mentor, mentee, mentorUUID, menteeUUID string) outputmodel.PairingOutput {
	return outputmodel.PairingOutput{
		ID:     id,
		UUID:   uuid.New().String(),
		Mentor: outputmodel.UserOutput{UUID: mentorUUID, Name: mentor},
		Mentee: outputmodel.UserOutput{UUID: menteeUUID, Name: mentee},
	}
}

func TestSortPlanningPairings_ByMentorThenMenteeThenID(t *testing.T) {
	list := []outputmodel.PairingOutput{
		pairingOut(1, "alpha", "m1", "m1", "e1"),
		pairingOut(2, "beta", "m2", "m2", "e2"),
		pairingOut(3, "beta", "m0", "m3", "e3"),
	}
	sortPlanningPairings(list)
	require.Equal(t, "alpha", list[0].Mentor.Name)
	require.Equal(t, "beta", list[1].Mentor.Name)
	require.Equal(t, "m0", list[1].Mentee.Name)
	require.Equal(t, uint(3), list[1].ID)
	require.Equal(t, uint(2), list[2].ID)
}

func TestFilterPlanningPairings_MentorMenteeSets(t *testing.T) {
	mid := uuid.New()
	eid := uuid.New()
	list := []outputmodel.PairingOutput{
		pairingOut(1, "a", "b", mid.String(), "mentee-1"),
		pairingOut(2, "c", "d", "mentor-2", eid.String()),
	}
	fMentor := PlanningPairingsListFilter{MentorUUIDs: []uuid.UUID{mid}}
	outM := filterPlanningPairings(list, fMentor, nil, nil)
	require.Len(t, outM, 1)
	require.Equal(t, uint(1), outM[0].ID)

	fMentee := PlanningPairingsListFilter{MenteeUUIDs: []uuid.UUID{eid}}
	outE := filterPlanningPairings(list, fMentee, nil, nil)
	require.Len(t, outE, 1)
	require.Equal(t, uint(2), outE[0].ID)
}

func TestFilterPlanningPairings_DepartmentOrTeam(t *testing.T) {
	list := []outputmodel.PairingOutput{
		pairingOut(1, "a", "b", "u1", "u2"),
		pairingOut(2, "c", "d", "u3", "u4"),
	}
	deptUsers := map[string]struct{}{"u1": {}}
	teamUsers := map[string]struct{}{"u4": {}}
	f := PlanningPairingsListFilter{DepartmentUUIDs: []uuid.UUID{uuid.New()}}
	out := filterPlanningPairings(list, f, deptUsers, nil)
	require.Len(t, out, 1)
	require.Equal(t, uint(1), out[0].ID)

	f2 := PlanningPairingsListFilter{TeamUUIDs: []uuid.UUID{uuid.New()}}
	out2 := filterPlanningPairings(list, f2, nil, teamUsers)
	require.Len(t, out2, 1)
	require.Equal(t, uint(2), out2[0].ID)
}

func TestSlicePlanningPairingsPage(t *testing.T) {
	list := make([]outputmodel.PairingOutput, 25)
	for i := range list {
		list[i] = pairingOut(uint(i+1), "a", "b", "x", "y")
	}
	page := slicePlanningPairingsPage(list, 2, PlanningPairingsPageSize)
	require.Len(t, page, 5)
	require.Empty(t, slicePlanningPairingsPage(list, 99, PlanningPairingsPageSize))
}