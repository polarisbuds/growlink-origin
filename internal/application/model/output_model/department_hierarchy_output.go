package outputmodel

type DepartmentHierarchyOutput struct {
	ParentDepartment DepartmentOutput `json:"parentDepartment"`
	ChildDepartment  DepartmentOutput `json:"childDepartment"`
}
