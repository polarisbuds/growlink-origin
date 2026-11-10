package outputmodel

type DepartmentOutput struct {
	ID   uint   `json:"-"`
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
