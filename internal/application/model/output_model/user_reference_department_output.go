package outputmodel

type UserReferenceDepartmentOutput struct {
	ID                    uint             `json:"id"`
	UserID                uint             `json:"userId"`
	ReferableDepartmentID uint             `json:"referableDepartmentId"`
	User                  UserOutput       `json:"user"`
	ReferableDepartment   DepartmentOutput `json:"referableDepartment"`
}
