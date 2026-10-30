package repository_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"growlink/internal/infrastructure/repository"
	"growlink/internal/shared/testutils"
	"growlink/internal/shared/testutils/fixtures"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// minimal valid 1x1 PNG
var tinyPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
	0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00,
	0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
	0x42, 0x60, 0x82,
}

func TestUserRepository_GetByID_DoesNotAttachImageBase64(t *testing.T) {
	tx := testutils.Begin(t)
	dir := t.TempDir()
	repo := repository.NewUserRepository(context.Background(), tx, dir)

	want := fixtures.User(t, tx, "icon-b64")
	iconName := "test-icon.png"
	want.ImageURI = &iconName
	require.NoError(t, repo.UpdateByID(want))

	require.NoError(t, os.WriteFile(filepath.Join(dir, iconName), tinyPNG, 0o644))

	got, err := repo.GetByID(want.ID)
	require.NoError(t, err)
	require.NotNil(t, got.ImageURI)
	assert.Equal(t, iconName, *got.ImageURI)
	assert.Nil(t, got.ImageBase64)
}

func TestUserRepository_GetByUUID_AttachesImageBase64(t *testing.T) {
	tx := testutils.Begin(t)
	dir := t.TempDir()
	repo := repository.NewUserRepository(context.Background(), tx, dir)

	want := fixtures.User(t, tx, "icon-uuid")
	iconName := "test-icon.png"
	want.ImageURI = &iconName
	require.NoError(t, repo.UpdateByID(want))

	require.NoError(t, os.WriteFile(filepath.Join(dir, iconName), tinyPNG, 0o644))

	got, err := repo.GetByUUID(want.UUID)
	require.NoError(t, err)
	require.NotNil(t, got.ImageURI)
	assert.Equal(t, iconName, *got.ImageURI)
	require.NotNil(t, got.ImageBase64)
	assert.True(t, strings.HasPrefix(*got.ImageBase64, "data:image/png;base64,"))
}
