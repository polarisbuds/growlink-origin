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
