package outputmodel

type UserReferenceTeamOutput struct {
	ID              uint       `json:"id"`
	UserID          uint       `json:"userId"`
	ReferableTeamID uint       `json:"referableTeamId"`
	User            UserOutput `json:"user"`
	ReferableTeam   TeamOutput `json:"referableTeam"`
}
