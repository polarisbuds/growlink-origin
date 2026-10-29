package inputmodel

import "growlink/internal/domain/valueobject"

type UserInput struct {
	Name         string                `json:"name"`
	Email        string                `json:"email"`
	ImageURI     *string               `json:"imageUri,omitempty"`
	RoleID       uint                  `json:"roleId"`
	DepartmentID uint                  `json:"departmentId"`
	Authority    valueobject.Authority `json:"authority"`
}
