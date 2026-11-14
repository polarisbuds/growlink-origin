package api

import (
	"fmt"
	"strconv"
	"strings"

	"growlink/internal/application/usecase/planning_usecase"
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func splitCommaUUIDs(s string) ([]uuid.UUID, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	out := make([]uuid.UUID, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := uuid.Parse(p)
		if err != nil {
			return nil, fmt.Errorf("invalid uuid")
		}
		out = append(out, id)
	}
	return out, nil
}

func parsePlanningMeetingsListQuery(c *echo.Context) (planning_usecase.PlanningMeetingsListFilter, int, bool, error) {
	filter := planning_usecase.PlanningMeetingsListFilter{}

	dir := strings.TrimSpace(c.QueryParam("dir"))
	dirDesc := true
	switch dir {
	case "", "desc":
		dirDesc = true
	case "asc":
		dirDesc = false
	default:
		return filter, 0, false, fmt.Errorf("invalid dir")
	}

	page := 1
	if ps := strings.TrimSpace(c.QueryParam("page")); ps != "" {
		p, err := strconv.Atoi(ps)
		if err != nil || p < 1 {
			return filter, 0, false, fmt.Errorf("invalid page")
		}
		page = p
	}

	var err error
	filter.MentorUUIDs, err = splitCommaUUIDs(c.QueryParam("mentorUuids"))
	if err != nil {
		return filter, 0, false, err
	}
	filter.MenteeUUIDs, err = splitCommaUUIDs(c.QueryParam("menteeUuids"))
	if err != nil {
		return filter, 0, false, err
	}
	filter.DepartmentUUIDs, err = splitCommaUUIDs(c.QueryParam("departmentUuids"))
	if err != nil {
		return filter, 0, false, err
	}

	if st := strings.TrimSpace(c.QueryParam("status")); st != "" {
		n, err := strconv.Atoi(st)
		if err != nil || n < 1 || n > 3 {
			return filter, 0, false, fmt.Errorf("invalid status")
		}
		status := valueobject.MeetingStatus(n)
		filter.Status = &status
	}

	fromStr := strings.TrimSpace(c.QueryParam("from"))
	toStr := strings.TrimSpace(c.QueryParam("to"))
	if fromStr != "" {
		from, err := planning_usecase.JSTDayStart(fromStr)
		if err != nil {
			return filter, 0, false, fmt.Errorf("invalid from")
		}
		filter.From = &from
	}
	if toStr != "" {
		to, err := planning_usecase.JSTDayEnd(toStr)
		if err != nil {
			return filter, 0, false, fmt.Errorf("invalid to")
		}
		filter.To = &to
	}
	if filter.From != nil && filter.To != nil && filter.From.After(*filter.To) {
		filter.DateRangeEmpty = true
	}

	return filter, page, dirDesc, nil
}