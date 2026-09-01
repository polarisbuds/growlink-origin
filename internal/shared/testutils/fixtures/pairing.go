package fixtures

import (
	"testing"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// createPairing inserts a Pairing with mentor/mentee users (Layer2).
func createPairing(t *testing.T, tx *gorm.DB) *entity.Pairing {
	t.Helper()

	mentor := createUser(t, tx, "mentor")
	mentee := createUser(t, tx, "mentee")

	pairing := &entity.Pairing{
		MentorID: mentor.ID,
		MenteeID: mentee.ID,
	}
	if err := tx.Create(pairing).Error; err != nil {
		t.Fatalf("failed to create pairing fixture: %v", err)
	}
	pairing.Mentor = mentor
	pairing.Mentee = mentee
	return pairing
}
