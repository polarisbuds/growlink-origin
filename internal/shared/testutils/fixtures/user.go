package fixtures

import (
	"fmt"
	"testing"
	"time"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// User inserts a User with Layer0 master dependencies (Layer1).
func User(t *testing.T, tx *gorm.DB, prefix string) *entity.User {
	t.Helper()
	return createUser(t, tx, prefix)
}

// createUser inserts a User with Layer0 master dependencies (Layer1).
func createUser(t *testing.T, tx *gorm.DB, prefix string) *entity.User {
	t.Helper()

	role := createRole(t, tx, prefix)
	department := createDepartment(t, tx, prefix)
	authority := createAuthority(t, tx, prefix)

	user := &entity.User{
		Name:         fmt.Sprintf("%s-name", prefix),
		Email:        fmt.Sprintf("%s-%d@example.com", prefix, time.Now().UnixNano()),
		RoleID:       role.ID,
		DepartmentID: department.ID,
		AuthorityID:  authority.ID,
	}
	if err := tx.Create(user).Error; err != nil {
		t.Fatalf("failed to create user fixture: %v", err)
	}
	return user
}
