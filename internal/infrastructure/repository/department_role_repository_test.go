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

func newDepartmentRoleRepo(t *testing.T, tx *gorm.DB) *repository.DepartmentRoleRepository {
	t.Helper()
	return repository.NewDepartmentRoleRepository(context.Background(), tx)
}

func TestDepartmentRoleRepository_Create(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newDepartmentRoleRepo(t, tx)

	user := fixtures.User(t, tx, "dept-role-create")
	department := fixtures.Team(t, tx).Department

	err := repo.Create(user.ID, department.ID, "Lead")
	assertions.NoError(t, err)

	got, err := repo.GetByUserID(user.ID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Lead", got[0].Name)
}

func TestDepartmentRoleRepository_Create_ForeignKey(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newDepartmentRoleRepo(t, tx)

	err := repo.Create(9_999_999, 9_999_998, "Lead")
	assertions.ForeignKeyViolation(t, err)
}

func TestDepartmentRoleRepository_Update(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newDepartmentRoleRepo(t, tx)

	want := fixtures.DepartmentRole(t, tx)
	want.Name = "Director"

	err := repo.Update(want.ID, want)
	assertions.NoError(t, err)

	got, err := repo.GetByDepartmentID(want.DepartmentID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Director", got[0].Name)
}

func TestDepartmentRoleRepository_Delete(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newDepartmentRoleRepo(t, tx)

	want := fixtures.DepartmentRole(t, tx)

	err := repo.Delete(want.ID)
	assertions.NoError(t, err)

	got, err := repo.GetByUserID(want.UserID)
	assertions.NoError(t, err)
	assert.Empty(t, got)
}

func TestDepartmentRoleRepository_GetAll(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newDepartmentRoleRepo(t, tx)

	fixtures.DepartmentRole(t, tx)

	got, err := repo.GetAll()
	assertions.NoError(t, err)
	require.NotEmpty(t, got)
	assert.NotNil(t, got[0].User)
	assert.NotNil(t, got[0].Department)
}
