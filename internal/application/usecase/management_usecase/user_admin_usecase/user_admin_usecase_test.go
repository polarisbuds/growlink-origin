package user_admin_usecase

import (
	"testing"

	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
)

func TestUserColumnChanges_Unchanged(t *testing.T) {
	suspended := false
	user := &entity.User{
		Authority:    valueobject.AuthorityGeneral,
		SystemRole:   valueobject.SystemRoleNone,
		DepartmentID: 1,
		RoleID:       2,
		Suspended:    false,
	}
	departments := []DepartmentAssignmentInput{{DepartmentID: 1, RoleID: 2}}
	columns, err := userColumnChanges(user, UpdateUserInput{
		Authority:   "general",
		SystemRole:  "none",
		Departments: departments,
		Suspended:   &suspended,
	}, departments)
	if err != nil {
		t.Fatal(err)
	}
	if len(columns) != 0 {
		t.Fatalf("expected no columns, got %#v", columns)
	}
}

func TestUserColumnChanges_OnlyChanged(t *testing.T) {
	suspended := false
	user := &entity.User{
		Authority:    valueobject.AuthorityGeneral,
		SystemRole:   valueobject.SystemRoleNone,
		DepartmentID: 1,
		RoleID:       2,
		Suspended:    true,
	}
	departments := []DepartmentAssignmentInput{{DepartmentID: 4, RoleID: 2}}
	columns, err := userColumnChanges(user, UpdateUserInput{
		Authority:   "boss",
		SystemRole:  "none",
		Suspended:   &suspended,
	}, departments)
	if err != nil {
		t.Fatal(err)
	}
	if columns["authority"] != valueobject.AuthorityBoss.Int() {
		t.Fatalf("authority = %#v", columns["authority"])
	}
	if columns["department_id"] != uint(4) {
		t.Fatalf("department_id = %#v", columns["department_id"])
	}
	if columns["suspended"] != false {
		t.Fatalf("suspended = %#v", columns["suspended"])
	}
	if _, ok := columns["system_role"]; ok {
		t.Fatal("unchanged system role was included")
	}
	if _, ok := columns["role_id"]; ok {
		t.Fatal("unchanged role was included")
	}
}

func TestDiffIDs_RemovedAndAdded(t *testing.T) {
	remove, add := diffIDs([]uint{1, 2}, []uint{2, 3})
	if !sameIDs(remove, []uint{1}) {
		t.Fatalf("remove = %#v", remove)
	}
	if !sameIDs(add, []uint{3}) {
		t.Fatalf("add = %#v", add)
	}
}

func TestDiffIDs_Unchanged(t *testing.T) {
	remove, add := diffIDs([]uint{2, 1}, []uint{1, 2})
	if len(remove) != 0 || len(add) != 0 {
		t.Fatalf("remove=%#v add=%#v", remove, add)
	}
}

func sameIDs(got, want []uint) bool {
	if len(got) != len(want) {
		return false
	}
	have := map[uint]int{}
	for _, id := range got {
		have[id]++
	}
	for _, id := range want {
		if have[id] == 0 {
			return false
		}
		have[id]--
	}
	return true
}

func TestValidateLineWorksLocalID(t *testing.T) {
	if err := ValidateLineWorksLocalID("user.name"); err != nil {
		t.Fatalf("expected valid id: %v", err)
	}
	if err := ValidateLineWorksLocalID("bad@id"); err == nil {
		t.Fatal("expected error for @ in id")
	}
}
