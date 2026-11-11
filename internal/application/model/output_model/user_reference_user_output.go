package outputmodel

type UserReferenceUserOutput struct {
	User          UserOutput `json:"user"`
	ReferableUser UserOutput `json:"referableUser"`
}
