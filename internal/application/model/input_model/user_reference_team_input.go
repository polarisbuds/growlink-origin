package inputmodel

type UserReferenceTeamInput struct {
	UserID          uint `json:"userId"`
	ReferableTeamID uint `json:"referableTeamId"`
}
