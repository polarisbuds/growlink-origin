package planning_usecase

import (
	"context"
	"slices"

	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"

	"github.com/google/uuid"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

const PlanningPairingsPageSize = 20

// PlanningPairingsListFilter is applied before paging or full list export.
type PlanningPairingsListFilter struct {
	MentorUUIDs     []uuid.UUID
	MenteeUUIDs     []uuid.UUID
	DepartmentUUIDs []uuid.UUID
	TeamUUIDs       []uuid.UUID
}

func filterPlanningPairings(
	list []outputmodel.PairingOutput,
	f PlanningPairingsListFilter,
	deptMemberUserUUIDs map[string]struct{},
	teamMemberUserUUIDs map[string]struct{},
) []outputmodel.PairingOutput {
	mentorSet := uuidSet(f.MentorUUIDs)
	menteeSet := uuidSet(f.MenteeUUIDs)
	deptFilter := len(f.DepartmentUUIDs) > 0
	teamFilter := len(f.TeamUUIDs) > 0

	out := make([]outputmodel.PairingOutput, 0, len(list))
	for _, p := range list {
		if len(mentorSet) > 0 {
			mu := p.Mentor.UUID
			if mu == "" || !mentorSet[mu] {
				continue
			}
		}
		if len(menteeSet) > 0 {
			mu := p.Mentee.UUID
			if mu == "" || !menteeSet[mu] {
				continue
			}
		}
		if deptFilter {
			mentorUUID := p.Mentor.UUID
			menteeUUID := p.Mentee.UUID
			_, okM := deptMemberUserUUIDs[mentorUUID]
			_, okE := deptMemberUserUUIDs[menteeUUID]
			if !okM && !okE {
				continue
			}
		}
		if teamFilter {
			mentorUUID := p.Mentor.UUID
			menteeUUID := p.Mentee.UUID
			_, okM := teamMemberUserUUIDs[mentorUUID]
			_, okE := teamMemberUserUUIDs[menteeUUID]
			if !okM && !okE {
				continue
			}
		}
		out = append(out, p)
	}
	return out
}

func sortPlanningPairings(list []outputmodel.PairingOutput) {
	coll := collate.New(language.Japanese)
	slices.SortFunc(list, func(a, b outputmodel.PairingOutput) int {
		if c := coll.CompareString(a.Mentor.Name, b.Mentor.Name); c != 0 {
			return c
		}
		if c := coll.CompareString(a.Mentee.Name, b.Mentee.Name); c != 0 {
			return c
		}
		if a.ID > b.ID {
			return -1
		}
		if a.ID < b.ID {
			return 1
		}
		return 0
	})
}

func slicePlanningPairingsPage(list []outputmodel.PairingOutput, page, pageSize int) []outputmodel.PairingOutput {
	if page < 1 || pageSize < 1 {
		return []outputmodel.PairingOutput{}
	}
	start := (page - 1) * pageSize
	if start >= len(list) {
		return []outputmodel.PairingOutput{}
	}
	end := start + pageSize
	if end > len(list) {
		end = len(list)
	}
	return slices.Clone(list[start:end])
}

func (uc *PlanningAccessUsecase) teamMemberUserUUIDsForTeams(teamUUIDs []string) (map[string]struct{}, error) {
	set := make(map[string]struct{})
	if uc.teamRepo == nil || len(teamUUIDs) == 0 {
		return set, nil
	}
	teams, err := uc.teamRepo.GetByUUIDs(teamUUIDs)
	if err != nil {
		return set, err
	}
	for _, team := range teams {
		if team.TeamMembers == nil {
			continue
		}
		for _, mem := range team.TeamMembers {
			if mem.Member != nil {
				set[mem.Member.UUID.String()] = struct{}{}
			}
		}
	}
	return set, nil
}

func (uc *PlanningAccessUsecase) listFilteredSortedPairings(
	userID uint,
	authority valueobject.Authority,
	systemRole valueobject.SystemRole,
	filter PlanningPairingsListFilter,
) ([]outputmodel.PairingOutput, error) {
	all, err := uc.ListPairingsForManagement(userID, authority, systemRole)
	if err != nil {
		return nil, err
	}
	deptUUIDs := []string{}
	for _, u := range filter.DepartmentUUIDs {
		deptUUIDs = append(deptUUIDs, u.String())
	}
	deptUsers, err := uc.deptMemberUserUUIDsForDepartments(deptUUIDs)
	if err != nil {
		return nil, err
	}
	teamUUIDs := []string{}
	for _, u := range filter.TeamUUIDs {
		teamUUIDs = append(teamUUIDs, u.String())
	}
	teamUsers, err := uc.teamMemberUserUUIDsForTeams(teamUUIDs)
	if err != nil {
		return nil, err
	}
	filtered := filterPlanningPairings(all, filter, deptUsers, teamUsers)
	sortPlanningPairings(filtered)
	return filtered, nil
}

func (uc *PlanningAccessUsecase) ListPairingsForManagementFiltered(
	userID uint,
	authority valueobject.Authority,
	systemRole valueobject.SystemRole,
	filter PlanningPairingsListFilter,
) (ret0 []outputmodel.PairingOutput, err error) {
	const op = "planning_usecase.ListPairingsForManagementFiltered"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	return uc.listFilteredSortedPairings(userID, authority, systemRole, filter)
}

func (uc *PlanningAccessUsecase) ListPairingsForManagementPaged(
	userID uint,
	authority valueobject.Authority,
	systemRole valueobject.SystemRole,
	filter PlanningPairingsListFilter,
	page int,
) (ret0 outputmodel.PlanningPairingsPageOutput, err error) {
	const op = "planning_usecase.ListPairingsForManagementPaged"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	filtered, err := uc.listFilteredSortedPairings(userID, authority, systemRole, filter)
	if err != nil {
		return outputmodel.PlanningPairingsPageOutput{}, err
	}
	items := slicePlanningPairingsPage(filtered, page, PlanningPairingsPageSize)
	return outputmodel.PlanningPairingsPageOutput{
		Items:    items,
		Total:    len(filtered),
		Page:     page,
		PageSize: PlanningPairingsPageSize,
	}, nil
}
