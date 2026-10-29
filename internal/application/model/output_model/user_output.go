package outputmodel

type UserOutput struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	ImageURI *string `json:"imageUri,omitempty"`
}
