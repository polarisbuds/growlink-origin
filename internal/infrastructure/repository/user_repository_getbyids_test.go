package repository_test

import (
	"testing"

	"growlink/internal/shared/testutils"
	"growlink/internal/shared/testutils/assertions"
	"growlink/internal/shared/testutils/fixtures"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_GetByIDs(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserRepo(t, tx)

	first := fixtures.User(t, tx, "ids-first")
	second := fixtures.User(t, tx, "ids-second")

	got, err := repo.GetByIDs([]uint{second.ID, first.ID})
	assertions.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, first.ID, got[0].ID)
	assert.Equal(t, second.ID, got[1].ID)
	assert.Equal(t, first.Name, got[0].Name)
	assert.Equal(t, second.Name, got[1].Name)
	assert.Equal(t, first.Email, got[0].Email)
	assert.Equal(t, second.Email, got[1].Email)
}

func TestUserRepository_GetByIDs_Empty(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserRepo(t, tx)

	got, err := repo.GetByIDs([]uint{})
	assertions.NoError(t, err)
	assert.Empty(t, got)
}

func TestUserRepository_GetByIDs_Partial(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserRepo(t, tx)

	want := fixtures.User(t, tx, "ids-partial")

	got, err := repo.GetByIDs([]uint{want.ID, 9_999_999})
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, want.ID, got[0].ID)
	assert.Equal(t, want.Name, got[0].Name)
}
