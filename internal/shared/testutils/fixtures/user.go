package fixtures

import (
	"fmt"
	"testing"
	"time"

	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"

	"gorm.io/gorm"
)

// User inserts a User with Layer0 master dependencies (Layer1).
func User(t *testing.T, tx *gorm.DB, prefix string) *entity.User {
	t.Helper()
	return createUser(t, tx, prefix)
}

// UserInDepartment inserts a User belonging to the given department.
func UserInDepartment(t *testing.T, tx *gorm.DB, prefix string, department *entity.Department) *entity.User {
	t.Helper()
	return createUserInDepartment(t, tx, prefix, department)
}

// createUser inserts a User with Layer0 master dependencies (Layer1).
func createUser(t *testing.T, tx *gorm.DB, prefix string) *entity.User {
	t.Helper()

	role := createRole(t, tx, prefix)
	department := createDepartment(t, tx, prefix)
	return createUserWithDeps(t, tx, prefix, role, department)
}

func createUserInDepartment(t *testing.T, tx *gorm.DB, prefix string, department *entity.Department) *entity.User {
	t.Helper()

	role := createRole(t, tx, prefix)
	return createUserWithDeps(t, tx, prefix, role, department)
}

func createUserWithDeps(t *testing.T, tx *gorm.DB, prefix string, role *entity.Role, department *entity.Department) *entity.User {
	t.Helper()

	user := &entity.User{
		Name:         fmt.Sprintf("%s-name", prefix),
		Email:        fmt.Sprintf("%s-%d@example.com", prefix, time.Now().UnixNano()),
		RoleID:       role.ID,
		DepartmentID: department.ID,
		Authority:    valueobject.AuthorityGeneral,
	}
	if err := tx.Create(user).Error; err != nil {
		t.Fatalf("failed to create user fixture: %v", err)
	}
	user.Role = role
	user.Department = department
	return user
}
