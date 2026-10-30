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

func newUserReferenceTeamRepo(t *testing.T, tx *gorm.DB) *repository.UserReferenceTeamRepository {
	t.Helper()
	return repository.NewUserReferenceTeamRepository(context.Background(), tx)
}

func TestUserReferenceTeamRepository_Create(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceTeamRepo(t, tx)

	user := fixtures.User(t, tx, "ref-team-create-owner")
	team := fixtures.Team(t, tx)

	err := repo.Create(user.ID, team.ID)
	assertions.NoError(t, err)

	got, err := repo.GetByUserID(user.ID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, team.ID, got[0].ReferableTeamID)
}

func TestUserReferenceTeamRepository_Create_ForeignKey(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceTeamRepo(t, tx)

	err := repo.Create(9_999_999, 9_999_998)
	assertions.ForeignKeyViolation(t, err)
}

func TestUserReferenceTeamRepository_Delete(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceTeamRepo(t, tx)

	want := fixtures.UserReferenceTeam(t, tx)

	err := repo.Delete(want.ID)
	assertions.NoError(t, err)

	got, err := repo.GetByUserID(want.UserID)
	assertions.NoError(t, err)
	assert.Empty(t, got)
}

func TestUserReferenceTeamRepository_GetByReferableTeamID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceTeamRepo(t, tx)

	want := fixtures.UserReferenceTeam(t, tx)

	got, err := repo.GetByReferableTeamID(want.ReferableTeamID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, want.ID, got[0].ID)
}

func TestUserReferenceTeamRepository_GetAll(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceTeamRepo(t, tx)

	fixtures.UserReferenceTeam(t, tx)

	got, err := repo.GetAll()
	assertions.NoError(t, err)
	require.NotEmpty(t, got)
	assert.NotNil(t, got[0].User)
	assert.NotNil(t, got[0].ReferableTeam)
}
