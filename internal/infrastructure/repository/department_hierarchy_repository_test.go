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

func newDepartmentHierarchyRepo(t *testing.T, tx *gorm.DB) *repository.DepartmentHierarchyRepository {
	t.Helper()
	return repository.NewDepartmentHierarchyRepository(context.Background(), tx)
}

func TestDepartmentHierarchyRepository_Create(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newDepartmentHierarchyRepo(t, tx)

	parentDept := fixtures.Team(t, tx).Department
	childDept := fixtures.Team(t, tx).Department

	err := repo.Create(parentDept.ID, childDept.ID)
	assertions.NoError(t, err)

	got, err := repo.GetByParentDepartmentID(parentDept.ID)
	assertions.NoError(t, err)
	require.NotEmpty(t, got)
	assert.Equal(t, parentDept.ID, got[0].ParentDepartmentID)
	assert.Equal(t, childDept.ID, got[0].ChildDepartmentID)
}

func TestDepartmentHierarchyRepository_Create_ForeignKey(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newDepartmentHierarchyRepo(t, tx)

	err := repo.Create(9_999_999, 9_999_998)
	assertions.ForeignKeyViolation(t, err)
}

func TestDepartmentHierarchyRepository_Delete(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newDepartmentHierarchyRepo(t, tx)

	want := fixtures.DepartmentHierarchy(t, tx)

	err := repo.Delete(want.ID)
	assertions.NoError(t, err)

	got, err := repo.GetByParentDepartmentID(want.ParentDepartmentID)
	assertions.NoError(t, err)
	assert.Empty(t, got)
}

func TestDepartmentHierarchyRepository_GetByChildDepartmentID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newDepartmentHierarchyRepo(t, tx)

	want := fixtures.DepartmentHierarchy(t, tx)

	got, err := repo.GetByChildDepartmentID(want.ChildDepartmentID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, want.ID, got[0].ID)
}

func TestDepartmentHierarchyRepository_GetAll(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newDepartmentHierarchyRepo(t, tx)

	fixtures.DepartmentHierarchy(t, tx)

	got, err := repo.GetAll()
	assertions.NoError(t, err)
	require.NotEmpty(t, got)
	assert.NotNil(t, got[0].ParentDepartment)
	assert.NotNil(t, got[0].ChildDepartment)
}
