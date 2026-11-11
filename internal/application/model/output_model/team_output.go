package outputmodel

type TeamOutput struct {
	Name        string             `json:"name"`
	Department  DepartmentOutput   `json:"department"`
	TeamMembers []TeamMemberOutput `json:"teamMembers"`
}
