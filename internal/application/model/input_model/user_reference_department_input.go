package inputmodel

type UserReferenceDepartmentInput struct {
	UserID                uint `json:"userId"`
	ReferableDepartmentID uint `json:"referableDepartmentId"`
}
