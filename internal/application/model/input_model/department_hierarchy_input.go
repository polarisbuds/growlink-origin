package inputmodel

type DepartmentHierarchyInput struct {
	ParentDepartmentID uint `json:"parentDepartmentId"`
	ChildDepartmentID  uint `json:"childDepartmentId"`
}
