package fixtures

import (
	"testing"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// DepartmentHierarchy inserts a DepartmentHierarchy with default dependencies.
func DepartmentHierarchy(t *testing.T, tx *gorm.DB) *entity.DepartmentHierarchy {
	t.Helper()
	return createDepartmentHierarchy(t, tx)
}

func createDepartmentHierarchy(t *testing.T, tx *gorm.DB) *entity.DepartmentHierarchy {
	t.Helper()

	parent := createDepartment(t, tx, "parent")
	child := createDepartment(t, tx, "child")

	hierarchy := &entity.DepartmentHierarchy{
		ParentDepartmentID: parent.ID,
		ChildDepartmentID:  child.ID,
	}
	if err := tx.Create(hierarchy).Error; err != nil {
		t.Fatalf("failed to create department hierarchy fixture: %v", err)
	}
	hierarchy.ParentDepartment = parent
	hierarchy.ChildDepartment = child
	return hierarchy
}
