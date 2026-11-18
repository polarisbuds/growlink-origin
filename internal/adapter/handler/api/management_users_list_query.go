package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type managementUsersListQuery struct {
	Name            string
	Email           string
	DepartmentUUIDs []uuid.UUID
	Page            int
	All             bool
	Size            int
}

func parseManagementUsersListQuery(c *echo.Context) (managementUsersListQuery, error) {
	q := managementUsersListQuery{Page: 1}
	q.Name = strings.TrimSpace(c.QueryParam("name"))
	q.Email = strings.TrimSpace(c.QueryParam("email"))
	q.All = strings.TrimSpace(c.QueryParam("all")) == "1"

	var err error
	q.DepartmentUUIDs, err = splitCommaUUIDs(c.QueryParam("departmentUuids"))
	if err != nil {
		return q, err
	}

	if ps := strings.TrimSpace(c.QueryParam("page")); ps != "" {
		p, err := strconv.Atoi(ps)
		if err != nil || p < 1 {
			return q, fmt.Errorf("invalid page")
		}
		q.Page = p
	}

	if size := strings.TrimSpace(c.QueryParam("pageSize")); size != "" {
		s, err := strconv.Atoi(size)
		if err != nil || s < 1 {
			return q, fmt.Errorf("invalid page size")
		}
		q.Size = s
	}

	return q, nil
}
