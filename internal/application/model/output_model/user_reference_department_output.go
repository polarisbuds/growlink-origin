package outputmodel

type UserReferenceDepartmentOutput struct {
	User                UserOutput       `json:"user"`
	ReferableDepartment DepartmentOutput `json:"referableDepartment"`
}
