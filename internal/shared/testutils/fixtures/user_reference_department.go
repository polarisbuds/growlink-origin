package fixtures

import (
	"testing"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// UserReferenceDepartment inserts a UserReferenceDepartment with default dependencies.
func UserReferenceDepartment(t *testing.T, tx *gorm.DB) *entity.UserReferenceDepartment {
	t.Helper()
	return createUserReferenceDepartment(t, tx)
}

func createUserReferenceDepartment(t *testing.T, tx *gorm.DB) *entity.UserReferenceDepartment {
	t.Helper()

	user := createUser(t, tx, "ref-dept-owner")
	department := createDepartment(t, tx, "ref-dept-target")

	reference := &entity.UserReferenceDepartment{
		UserID:                user.ID,
		ReferableDepartmentID: department.ID,
	}
	if err := tx.Create(reference).Error; err != nil {
		t.Fatalf("failed to create user reference department fixture: %v", err)
	}
	reference.User = user
	reference.ReferableDepartment = department
	return reference
}
