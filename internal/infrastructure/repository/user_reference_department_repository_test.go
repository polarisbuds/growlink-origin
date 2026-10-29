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

func newUserReferenceDepartmentRepo(t *testing.T, tx *gorm.DB) *repository.UserReferenceDepartmentRepository {
	t.Helper()
	return repository.NewUserReferenceDepartmentRepository(context.Background(), tx)
}

func TestUserReferenceDepartmentRepository_Create(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceDepartmentRepo(t, tx)

	user := fixtures.User(t, tx, "ref-dept-create-owner")
	department := fixtures.Team(t, tx).Department

	err := repo.Create(user.ID, department.ID)
	assertions.NoError(t, err)

	got, err := repo.GetByUserID(user.ID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, department.ID, got[0].ReferableDepartmentID)
}

func TestUserReferenceDepartmentRepository_Create_ForeignKey(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceDepartmentRepo(t, tx)

	err := repo.Create(9_999_999, 9_999_998)
	assertions.ForeignKeyViolation(t, err)
}

func TestUserReferenceDepartmentRepository_Delete(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceDepartmentRepo(t, tx)

	want := fixtures.UserReferenceDepartment(t, tx)

	err := repo.Delete(want.ID)
	assertions.NoError(t, err)

	got, err := repo.GetByUserID(want.UserID)
	assertions.NoError(t, err)
	assert.Empty(t, got)
}

func TestUserReferenceDepartmentRepository_GetByReferableDepartmentID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceDepartmentRepo(t, tx)

	want := fixtures.UserReferenceDepartment(t, tx)

	got, err := repo.GetByReferableDepartmentID(want.ReferableDepartmentID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, want.ID, got[0].ID)
}

func TestUserReferenceDepartmentRepository_GetAll(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserReferenceDepartmentRepo(t, tx)

	fixtures.UserReferenceDepartment(t, tx)

	got, err := repo.GetAll()
	assertions.NoError(t, err)
	require.NotEmpty(t, got)
	assert.NotNil(t, got[0].User)
	assert.NotNil(t, got[0].ReferableDepartment)
}
