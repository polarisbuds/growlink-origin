package outputmodel

type DepartmentRoleOutput struct {
	ID           uint             `json:"id"`
	Name         string           `json:"name"`
	UserID       uint             `json:"userId"`
	DepartmentID uint             `json:"departmentId"`
	User         UserOutput       `json:"user"`
	Department   DepartmentOutput `json:"department"`
}
