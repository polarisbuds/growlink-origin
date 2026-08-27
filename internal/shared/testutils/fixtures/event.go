package fixtures

import (
	"fmt"
	"testing"
	"time"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// Event inserts an Event master record.
func Event(t *testing.T, tx *gorm.DB, prefix string) *entity.Event {
	t.Helper()
	return createEvent(t, tx, prefix)
}

func createEvent(t *testing.T, tx *gorm.DB, prefix string) *entity.Event {
	t.Helper()

	event := &entity.Event{
		Name: fmt.Sprintf("%s-event-%d", prefix, time.Now().UnixNano()),
	}
	if err := tx.Create(event).Error; err != nil {
		t.Fatalf("failed to create event fixture: %v", err)
	}
	return event
}
