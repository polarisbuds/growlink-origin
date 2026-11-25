package planning_usecase

import (
	"context"
	"slices"
	"time"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"

	"github.com/google/uuid"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

const PlanningMeetingsPageSize = 20

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

// PlanningMeetingsListFilter is applied before paging.
type PlanningMeetingsListFilter struct {
	MentorUUIDs     []uuid.UUID
	MenteeUUIDs     []uuid.UUID
	DepartmentUUIDs []uuid.UUID
	Status          *valueobject.MeetingStatus
	From            *time.Time // start of day JST inclusive
	To              *time.Time // end of day JST inclusive
	DateRangeEmpty  bool       // from after to
	IncludeDeleted  bool       // include meetings whose pairing is soft-deleted
}

func JSTDayStart(dateYYYYMMDD string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", dateYYYYMMDD, jst)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func JSTDayEnd(dateYYYYMMDD string) (time.Time, error) {
	start, err := JSTDayStart(dateYYYYMMDD)
	if err != nil {
		return time.Time{}, err
	}
	return start.Add(24 * time.Hour).Add(-time.Nanosecond), nil
}

func planningMeetingsRepoFilter(f PlanningMeetingsListFilter) repository_interface.MeetingListFilter {
	mentors := make([]string, len(f.MentorUUIDs))
	for idx, id := range f.MentorUUIDs {
		mentors[idx] = id.String()
	}
	mentees := make([]string, len(f.MenteeUUIDs))
	for idx, id := range f.MenteeUUIDs {
		mentees[idx] = id.String()
	}
	depts := make([]string, len(f.DepartmentUUIDs))
	for idx, id := range f.DepartmentUUIDs {
		depts[idx] = id.String()
	}
	return repository_interface.MeetingListFilter{
		MentorUUIDs:     mentors,
		MenteeUUIDs:     mentees,
		DepartmentUUIDs: depts,
		Status:          f.Status,
		From:            f.From,
		To:              f.To,
		IncludeDeleted:  f.IncludeDeleted,
	}
}

func sortPlanningMeetings(list []outputmodel.MeetingOutput, desc bool) {
	coll := collate.New(language.Japanese)
	slices.SortFunc(list, func(a, b outputmodel.MeetingOutput) int {
		if c := compareImplementedAt(a.ImplementedAt, b.ImplementedAt, desc); c != 0 {
			return c
		}
		if c := coll.CompareString(a.Pairing.Mentor.Name, b.Pairing.Mentor.Name); c != 0 {
			return c
		}
		if c := coll.CompareString(a.Pairing.Mentee.Name, b.Pairing.Mentee.Name); c != 0 {
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

func compareImplementedAt(a, b time.Time, desc bool) int {
	aZero := a.IsZero()
	bZero := b.IsZero()
	if aZero && bZero {
		return 0
	}
	if aZero {
		if desc {
			return 1
		}
		return -1
	}
	if bZero {
		if desc {
			return -1
		}
		return 1
	}
	if desc {
		if a.After(b) {
			return -1
		}
		if a.Before(b) {
			return 1
		}
		return 0
	}
	if a.Before(b) {
		return -1
	}
	if a.After(b) {
		return 1
	}
	return 0
}

func slicePlanningMeetingsPage(list []outputmodel.MeetingOutput, page, pageSize int) []outputmodel.MeetingOutput {
	if page < 1 || pageSize < 1 {
		return []outputmodel.MeetingOutput{}
	}
	start := (page - 1) * pageSize
	if start >= len(list) {
		return []outputmodel.MeetingOutput{}
	}
	end := start + pageSize
	if end > len(list) {
		end = len(list)
	}
	return slices.Clone(list[start:end])
}

func (uc *PlanningAccessUsecase) deptMemberUserUUIDsForDepartments(deptUUIDs []string) (map[string]struct{}, error) {
	set := make(map[string]struct{})
	if uc.deptRepo == nil || len(deptUUIDs) == 0 {
		return set, nil
	}
	depts, err := uc.deptRepo.GetByUUIDs(deptUUIDs)
	if err != nil {
		return set, err
	}

	for _, dept := range depts {
		for _, mem := range dept.DepartmentMembers {
			if mem.User != nil {
				set[mem.User.UUID.String()] = struct{}{}
			}
		}
	}
	return set, nil
}

func (uc *PlanningAccessUsecase) ListPlanningMeetingsPaged(
	userID uint,
	authority valueobject.Authority,
	systemRole valueobject.SystemRole,
	filter PlanningMeetingsListFilter,
	page int,
	dirDesc bool,
) (ret0 outputmodel.PlanningMeetingsPageOutput, err error) {
	const op = "planning_usecase.ListPlanningMeetingsPaged"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()

	if filter.DateRangeEmpty {
		return outputmodel.PlanningMeetingsPageOutput{
			Items:    []outputmodel.MeetingOutput{},
			Total:    0,
			Page:     page,
			PageSize: PlanningMeetingsPageSize,
		}, nil
	}

	repoFilter := planningMeetingsRepoFilter(filter)
	if fullPlanningAccess(systemRole) || authority.IsOwner() {
		repoFilter.LimitToMeetingIDs = false
	} else {
		scoped, err := uc.ListPlanningMeetings(userID, authority, systemRole)
		if err != nil {
			return outputmodel.PlanningMeetingsPageOutput{}, err
		}
		repoFilter.MeetingIDs = make([]uint, len(scoped))
		for idx, m := range scoped {
			repoFilter.MeetingIDs[idx] = m.ID
		}
		repoFilter.LimitToMeetingIDs = true
	}

	meetings, err := uc.meetingRepo.GetByFilter(repoFilter)
	if err != nil {
		return outputmodel.PlanningMeetingsPageOutput{}, err
	}
	meetingOutputs := converter.MeetingEntitiesToOutputs(meetings)
	sortPlanningMeetings(meetingOutputs, dirDesc)
	items := slicePlanningMeetingsPage(meetingOutputs, page, PlanningMeetingsPageSize)
	return outputmodel.PlanningMeetingsPageOutput{
		Items:    items,
		Total:    len(meetingOutputs),
		Page:     page,
		PageSize: PlanningMeetingsPageSize,
	}, nil
}