package outputmodel

type DepartmentHierarchyOutput struct {
	ID                 uint             `json:"id"`
	ParentDepartmentID uint             `json:"parentDepartmentId"`
	ChildDepartmentID  uint             `json:"childDepartmentId"`
	ParentDepartment   DepartmentOutput `json:"parentDepartment"`
	ChildDepartment    DepartmentOutput `json:"childDepartment"`
}
