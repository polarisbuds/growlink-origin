package fixtures

import (
	"testing"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// createPeering inserts a Peering with mentor/mentee users (Layer2).
func createPeering(t *testing.T, tx *gorm.DB) *entity.Peering {
	t.Helper()

	mentor := createUser(t, tx, "mentor")
	mentee := createUser(t, tx, "mentee")

	peering := &entity.Peering{
		MentorID: int64(mentor.ID),
		MenteeID: int64(mentee.ID),
	}
	if err := tx.Create(peering).Error; err != nil {
		t.Fatalf("failed to create peering fixture: %v", err)
	}
	peering.Mentor = mentor
	peering.Mentee = mentee
	return peering
}
