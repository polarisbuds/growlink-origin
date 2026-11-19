package repository_test

import (
	"context"
	"os"
	"testing"

	"growlink/internal/shared/testutils/containers"
)

func TestMain(m *testing.M) {
	code := m.Run()

	ctx := context.Background()
	if err := containers.Terminate(ctx); err != nil {
		os.Exit(1)
	}

	os.Exit(code)
}
