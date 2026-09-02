package fixtures

import (
	"testing"

	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"

	"gorm.io/gorm"
)

// History inserts a History with default dependencies.
func History(t *testing.T, tx *gorm.DB) *entity.History {
	t.Helper()
	return createHistory(t, tx)
}

func createHistory(t *testing.T, tx *gorm.DB) *entity.History {
	t.Helper()

	user := createUser(t, tx, "history")

	history := &entity.History{
		UserID:    user.ID,
		EventType: valueobject.EventTypeLogin,
	}
	if err := tx.Create(history).Error; err != nil {
		t.Fatalf("failed to create history fixture: %v", err)
	}
	history.User = user
	return history
}
