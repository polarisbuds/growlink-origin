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

func newUserReferenceUserRepo(t *testing.T, tx *gorm.DB) *repository.UserReferenceUserRepository {
	t.Helper()
	return repository.NewUserReferenceUserRepository(context.Background(), tx)
}

func TestUserReferenceUserRepository_Create(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceUserRepo(t, tx)

	user := fixtures.User(t, tx, "ref-user-create-owner")
	referableUser := fixtures.User(t, tx, "ref-user-create-target")

	err := repo.Create(user.ID, referableUser.ID)
	assertions.NoError(t, err)

	got, err := repo.GetByUserID(user.ID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, referableUser.ID, got[0].ReferableUserID)
}

func TestUserReferenceUserRepository_Create_ForeignKey(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceUserRepo(t, tx)

	err := repo.Create(9_999_999, 9_999_998)
	assertions.ForeignKeyViolation(t, err)
}

func TestUserReferenceUserRepository_Delete(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceUserRepo(t, tx)

	want := fixtures.UserReferenceUser(t, tx)

	err := repo.Delete(want.ID)
	assertions.NoError(t, err)

	got, err := repo.GetByUserID(want.UserID)
	assertions.NoError(t, err)
	assert.Empty(t, got)
}

func TestUserReferenceUserRepository_GetByReferableUserID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceUserRepo(t, tx)

	want := fixtures.UserReferenceUser(t, tx)

	got, err := repo.GetByReferableUserID(want.ReferableUserID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, want.ID, got[0].ID)
}

func TestUserReferenceUserRepository_GetAll(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceUserRepo(t, tx)

	fixtures.UserReferenceUser(t, tx)

	got, err := repo.GetAll()
	assertions.NoError(t, err)
	require.NotEmpty(t, got)
	assert.NotNil(t, got[0].User)
	assert.NotNil(t, got[0].ReferableUser)
}
