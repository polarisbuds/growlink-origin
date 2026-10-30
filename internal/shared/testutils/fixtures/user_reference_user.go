package fixtures

import (
	"testing"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// UserReferenceUser inserts a UserReferenceUser with default dependencies.
func UserReferenceUser(t *testing.T, tx *gorm.DB) *entity.UserReferenceUser {
	t.Helper()
	return createUserReferenceUser(t, tx)
}

func createUserReferenceUser(t *testing.T, tx *gorm.DB) *entity.UserReferenceUser {
	t.Helper()

	user := createUser(t, tx, "ref-user-owner")
	referableUser := createUser(t, tx, "ref-user-target")

	reference := &entity.UserReferenceUser{
		UserID:          user.ID,
		ReferableUserID: referableUser.ID,
	}
	if err := tx.Create(reference).Error; err != nil {
		t.Fatalf("failed to create user reference user fixture: %v", err)
	}
	reference.User = user
	reference.ReferableUser = referableUser
	return reference
}
