package outputmodel

type TeamMemberOutput struct {
	ID       uint       `json:"id"`
	IsReader bool       `json:"isReader"`
	TeamID   uint       `json:"teamId"`
	MemberID uint       `json:"memberId"`
	Team     TeamOutput `json:"team"`
	Member   UserOutput `json:"member"`
}
