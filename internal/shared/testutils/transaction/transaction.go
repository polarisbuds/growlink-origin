package transaction

import (
	"testing"

	"gorm.io/gorm"
)

// Begin starts a transaction and registers rollback via t.Cleanup.
func Begin(t *testing.T, db *gorm.DB) *gorm.DB {
	t.Helper()

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	t.Cleanup(func() {
		if err := tx.Rollback().Error; err != nil {
			t.Fatalf("failed to rollback transaction: %v", err)
		}
	})

	return tx
}
