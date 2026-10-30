package outputmodel

type TeamOutput struct {
	ID         uint               `json:"id"`
	Name       string             `json:"name"`
	DepartmentID uint             `json:"departmentId"`
	Department DepartmentOutput   `json:"department"`
	TeamMembers []TeamMemberOutput `json:"teamMembers"`
}
