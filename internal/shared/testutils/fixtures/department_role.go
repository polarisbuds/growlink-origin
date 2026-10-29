package fixtures

import (
	"fmt"
	"testing"
	"time"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// DepartmentRole inserts a DepartmentRole with default dependencies.
func DepartmentRole(t *testing.T, tx *gorm.DB) *entity.DepartmentRole {
	t.Helper()
	return createDepartmentRole(t, tx, "manager")
}

func createDepartmentRole(t *testing.T, tx *gorm.DB, name string) *entity.DepartmentRole {
	t.Helper()

	user := createUser(t, tx, "dept-role")
	department := createDepartment(t, tx, "dept-role")

	role := &entity.DepartmentRole{
		Name:         fmt.Sprintf("%s-%d", name, time.Now().UnixNano()),
		UserID:       user.ID,
		DepartmentID: department.ID,
	}
	if err := tx.Create(role).Error; err != nil {
		t.Fatalf("failed to create department role fixture: %v", err)
	}
	role.User = user
	role.Department = department
	return role
}
