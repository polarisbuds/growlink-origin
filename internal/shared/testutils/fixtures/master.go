package fixtures

import (
	"fmt"
	"testing"
	"time"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// createRole inserts a Role master record (Layer0).
func createRole(t *testing.T, tx *gorm.DB, prefix string) *entity.Role {
	t.Helper()

	role := &entity.Role{Name: fmt.Sprintf("%s-role-%d", prefix, time.Now().UnixNano())}
	if err := tx.Create(role).Error; err != nil {
		t.Fatalf("failed to create role fixture: %v", err)
	}
	return role
}

// createDepartment inserts a Department master record (Layer0).
func createDepartment(t *testing.T, tx *gorm.DB, prefix string) *entity.Department {
	t.Helper()

	department := &entity.Department{Name: fmt.Sprintf("%s-dept-%d", prefix, time.Now().UnixNano())}
	if err := tx.Create(department).Error; err != nil {
		t.Fatalf("failed to create department fixture: %v", err)
	}
	return department
}
