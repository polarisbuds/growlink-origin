package planning_usecase

import (
	"slices"
	"time"

	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/valueobject"

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
	return start.Add(24*time.Hour).Add(-time.Nanosecond), nil
}

func filterPlanningMeetings(
	list []outputmodel.MeetingOutput,
	f PlanningMeetingsListFilter,
	deptMemberUserUUIDs map[string]struct{},
) []outputmodel.MeetingOutput {
	if f.DateRangeEmpty {
		return []outputmodel.MeetingOutput{}
	}
	hasDateFilter := f.From != nil || f.To != nil
	mentorSet := uuidSet(f.MentorUUIDs)
	menteeSet := uuidSet(f.MenteeUUIDs)
	deptFilter := len(f.DepartmentUUIDs) > 0

	out := make([]outputmodel.MeetingOutput, 0, len(list))
	for _, m := range list {
		if len(mentorSet) > 0 {
			mu := m.Pairing.Mentor.UUID
			if mu == "" || !mentorSet[mu] {
				continue
			}
		}
		if len(menteeSet) > 0 {
			mu := m.Pairing.Mentee.UUID
			if mu == "" || !menteeSet[mu] {
				continue
			}
		}
		if deptFilter {
			mentorUUID := m.Pairing.Mentor.UUID
			menteeUUID := m.Pairing.Mentee.UUID
			_, okM := deptMemberUserUUIDs[mentorUUID]
			_, okE := deptMemberUserUUIDs[menteeUUID]
			if !okM && !okE {
				continue
			}
		}
		if f.Status != nil && m.Status != *f.Status {
			continue
		}
		if hasDateFilter {
			if m.ImplementedAt.IsZero() {
				continue
			}
			at := m.ImplementedAt.In(jst)
			if f.From != nil && at.Before(*f.From) {
				continue
			}
			if f.To != nil && at.After(*f.To) {
				continue
			}
		}
		out = append(out, m)
	}
	return out
}

func uuidSet(ids []uuid.UUID) map[string]bool {
	if len(ids) == 0 {
		return nil
	}
	s := make(map[string]bool, len(ids))
	for _, id := range ids {
		s[id.String()] = true
	}
	return s
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

func (uc *PlanningAccessUsecase) deptMemberUserUUIDsForDepartments(deptUUIDs []uuid.UUID) (map[string]struct{}, error) {
	set := make(map[string]struct{})
	if uc.deptRepo == nil || len(deptUUIDs) == 0 {
		return set, nil
	}
	for _, du := range deptUUIDs {
		dept, err := uc.deptRepo.GetByUUID(du)
		if err != nil {
			continue
		}
		members, err := uc.deptMemberRepo.GetByDepartmentID(dept.ID)
		if err != nil {
			return nil, err
		}
		for _, mem := range members {
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
) (outputmodel.PlanningMeetingsPageOutput, error) {
	all, err := uc.ListPlanningMeetings(userID, authority, systemRole)
	if err != nil {
		return outputmodel.PlanningMeetingsPageOutput{}, err
	}
	deptUsers, err := uc.deptMemberUserUUIDsForDepartments(filter.DepartmentUUIDs)
	if err != nil {
		return outputmodel.PlanningMeetingsPageOutput{}, err
	}
	filtered := filterPlanningMeetings(all, filter, deptUsers)
	sortPlanningMeetings(filtered, dirDesc)
	items := slicePlanningMeetingsPage(filtered, page, PlanningMeetingsPageSize)
	return outputmodel.PlanningMeetingsPageOutput{
		Items:    items,
		Total:    len(filtered),
		Page:     page,
		PageSize: PlanningMeetingsPageSize,
	}, nil
}