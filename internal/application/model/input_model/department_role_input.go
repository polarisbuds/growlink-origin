package inputmodel

type DepartmentRoleInput struct {
	Name         string `json:"name"`
	UserID       uint   `json:"userId"`
	DepartmentID uint   `json:"departmentId"`
}
