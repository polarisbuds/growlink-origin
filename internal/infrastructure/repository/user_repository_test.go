package repository_test

import (
	"context"
	"testing"

	"growlink/internal/infrastructure/repository"
	"growlink/internal/shared/testutils"
	"growlink/internal/shared/testutils/assertions"
	"growlink/internal/shared/testutils/fixtures"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newUserRepo(t *testing.T, tx *gorm.DB) *repository.UserRepository {
	t.Helper()
	return repository.NewUserRepository(context.Background(), tx)
}

func TestUserRepository_GetByID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserRepo(t, tx)

	want := fixtures.User(t, tx, "get")

	got, err := repo.GetByID(want.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, want.ID, got.ID)
	assertions.EqualUUID(t, want.UUID, got.UUID)
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Email, got.Email)
	assert.Equal(t, want.RoleID, got.RoleID)
	assert.Equal(t, want.DepartmentID, got.DepartmentID)
	assert.Equal(t, want.AuthorityID, got.AuthorityID)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserRepo(t, tx)

	_, err := repo.GetByID(9_999_999)
	assertions.NotFound(t, err)
}

func TestUserRepository_UpdateByID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserRepo(t, tx)

	want := fixtures.User(t, tx, "update")
	imageURI := "icon.png"
	want.Name = "updated-name"
	want.ImageURI = &imageURI

	err := repo.UpdateByID(want)
	assertions.NoError(t, err)

	got, err := repo.GetByID(want.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, "updated-name", got.Name)
	require.NotNil(t, got.ImageURI)
	assert.Equal(t, imageURI, *got.ImageURI)
	assert.Equal(t, want.Email, got.Email)
}

func TestUserRepository_UpdateByID_Nil(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserRepo(t, tx)

	err := repo.UpdateByID(nil)
	assertions.NoError(t, err)
}
