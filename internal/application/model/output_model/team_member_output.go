package outputmodel

type TeamMemberOutput struct {
	IsReader bool       `json:"isReader"`
	Team     TeamOutput `json:"team"`
	Member   UserOutput `json:"member"`
}
