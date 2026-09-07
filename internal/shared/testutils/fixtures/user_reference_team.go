package fixtures

import (
	"testing"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// UserReferenceTeam inserts a UserReferenceTeam with default dependencies.
func UserReferenceTeam(t *testing.T, tx *gorm.DB) *entity.UserReferenceTeam {
	t.Helper()
	return createUserReferenceTeam(t, tx)
}

func createUserReferenceTeam(t *testing.T, tx *gorm.DB) *entity.UserReferenceTeam {
	t.Helper()

	user := createUser(t, tx, "ref-team-owner")
	team := createTeam(t, tx)

	reference := &entity.UserReferenceTeam{
		UserID:          user.ID,
		ReferableTeamID: team.ID,
	}
	if err := tx.Create(reference).Error; err != nil {
		t.Fatalf("failed to create user reference team fixture: %v", err)
	}
	reference.User = user
	reference.ReferableTeam = team
	return reference
}
