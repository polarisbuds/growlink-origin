package testutils

import (
	"context"
	"testing"

	"growlink/internal/shared/testutils/database"
	"growlink/internal/shared/testutils/transaction"

	"gorm.io/gorm"
)

// Begin starts a test transaction using the shared PostgreSQL container and GORM DB.
func Begin(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := database.DB(context.Background())
	if err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}

	return transaction.Begin(t, db)
}
