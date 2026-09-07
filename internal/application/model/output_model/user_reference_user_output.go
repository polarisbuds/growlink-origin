package outputmodel

type UserReferenceUserOutput struct {
	ID              uint       `json:"id"`
	UserID          uint       `json:"userId"`
	ReferableUserID uint       `json:"referableUserId"`
	User            UserOutput `json:"user"`
	ReferableUser   UserOutput `json:"referableUser"`
}
