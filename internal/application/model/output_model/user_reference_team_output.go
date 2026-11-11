package outputmodel

type UserReferenceTeamOutput struct {
	User          UserOutput `json:"user"`
	ReferableTeam TeamOutput `json:"referableTeam"`
}
