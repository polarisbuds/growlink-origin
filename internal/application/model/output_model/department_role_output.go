package outputmodel

type DepartmentRoleOutput struct {
	Name       string           `json:"name"`
	User       UserOutput       `json:"user"`
	Department DepartmentOutput `json:"department"`
}
