package user_admin_usecase

import (
	"growlink/internal/domain/entity"
	"strings"
)

const ManagementUsersPageSize = 20

func matchesManagementUserFilter(
	u entity.User,
	nameQ string,
	emailQ string,
	deptMemberUserIDs map[uint]struct{},
	hasDeptFilter bool,
) bool {
	if nameQ != "" && !strings.Contains(strings.ToLower(u.Name), strings.ToLower(nameQ)) {
		return false
	}
	if emailQ != "" && !strings.Contains(strings.ToLower(u.Email), strings.ToLower(emailQ)) {
		return false
	}
	if hasDeptFilter {
		if _, ok := deptMemberUserIDs[u.ID]; !ok {
			return false
		}
	}
	return true
}

func filterManagementUsers(
	users []entity.User,
	nameQ string,
	emailQ string,
	deptMemberUserIDs map[uint]struct{},
	hasDeptFilter bool,
) []entity.User {
	nameQ = strings.TrimSpace(nameQ)
	emailQ = strings.TrimSpace(emailQ)
	out := make([]entity.User, 0)
	for _, u := range users {
		if matchesManagementUserFilter(u, nameQ, emailQ, deptMemberUserIDs, hasDeptFilter) {
			out = append(out, u)
		}
	}
	return out
}

func sliceManagementUsersPage(list []entity.User, page, pageSize int) []entity.User {
	if page < 1 || pageSize < 1 {
		return []entity.User{}
	}
	start := (page - 1) * pageSize
	if start >= len(list) {
		return []entity.User{}
	}
	end := start + pageSize
	if end > len(list) {
		end = len(list)
	}
	return list[start:end]
}
