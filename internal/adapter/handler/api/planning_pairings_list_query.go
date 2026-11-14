package api

import (
	"fmt"
	"strconv"
	"strings"

	"growlink/internal/application/usecase/planning_usecase"

	"github.com/labstack/echo/v5"
)

func parsePlanningPairingsListFilter(c *echo.Context) (planning_usecase.PlanningPairingsListFilter, error) {
	filter := planning_usecase.PlanningPairingsListFilter{}
	var err error
	filter.MentorUUIDs, err = splitCommaUUIDs(c.QueryParam("mentorUuids"))
	if err != nil {
		return filter, err
	}
	filter.MenteeUUIDs, err = splitCommaUUIDs(c.QueryParam("menteeUuids"))
	if err != nil {
		return filter, err
	}
	filter.DepartmentUUIDs, err = splitCommaUUIDs(c.QueryParam("departmentUuids"))
	if err != nil {
		return filter, err
	}
	filter.TeamUUIDs, err = splitCommaUUIDs(c.QueryParam("teamUuids"))
	if err != nil {
		return filter, err
	}
	return filter, nil
}

func pairingsListQueryHasFilters(c *echo.Context) bool {
	for _, key := range []string{"mentorUuids", "menteeUuids", "departmentUuids", "teamUuids"} {
		if strings.TrimSpace(c.QueryParam(key)) != "" {
			return true
		}
	}
	return false
}

func parsePairingsListPage(c *echo.Context) (int, bool, error) {
	pageStr := strings.TrimSpace(c.QueryParam("page"))
	if pageStr == "" {
		return 0, false, nil
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, true, fmt.Errorf("invalid page")
	}
	return page, true, nil
}

func planningPairingsFilterEmpty(f planning_usecase.PlanningPairingsListFilter) bool {
	return len(f.MentorUUIDs) == 0 &&
		len(f.MenteeUUIDs) == 0 &&
		len(f.DepartmentUUIDs) == 0 &&
		len(f.TeamUUIDs) == 0
}