package inputmodel

type TeamMemberInput struct {
	IsReader bool `json:"isReader"`
	TeamID   uint `json:"teamId"`
	MemberID uint `json:"memberId"`
}
