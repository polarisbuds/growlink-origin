package user_admin_usecase

import (
	"testing"

	"growlink/internal/domain/entity"
)

func TestMatchesManagementUserFilter_NameEmailAndDept(t *testing.T) {
	deptUsers := map[uint]struct{}{2: {}}
	u := entity.User{Name: "Alice", Email: "a@x.com"}
	u.ID = 2
	if !matchesManagementUserFilter(u, "ali", "x.com", deptUsers, true) {
		t.Fatal("expected match")
	}
	if matchesManagementUserFilter(u, "bob", "", deptUsers, true) {
		t.Fatal("name mismatch")
	}
	u2 := entity.User{Name: "Alice", Email: "a@x.com"}
	u2.ID = 99
	if matchesManagementUserFilter(u2, "", "", deptUsers, true) {
		t.Fatal("user not in dept set")
	}
	if !matchesManagementUserFilter(u2, "", "", deptUsers, false) {
		t.Fatal("should match without dept filter")
	}
}

func TestSliceManagementUsersPage(t *testing.T) {
	users := make([]entity.User, 3)
	for i := range users {
		users[i].ID = uint(i + 1)
	}
	page := sliceManagementUsersPage(users, 2, 2)
	if len(page) != 1 || page[0].ID != 3 {
		t.Fatalf("page = %#v", page)
	}
	if len(sliceManagementUsersPage(users, 5, 2)) != 0 {
		t.Fatal("expected empty page")
	}
}

func TestUserIDsInDepartments(t *testing.T) {
	deptSet := map[string]struct{}{"d1": {}}
	deptIDToUUID := map[uint]string{10: "d1", 11: "d2"}
	members := []entity.DepartmentMember{
		{DepartmentID: 10, UserID: 5},
		{DepartmentID: 11, UserID: 6},
	}
	set := userIDsInDepartments(members, deptSet, deptIDToUUID)
	if len(set) != 1 || set[5] != struct{}{} {
		t.Fatalf("set = %#v", set)
	}
}