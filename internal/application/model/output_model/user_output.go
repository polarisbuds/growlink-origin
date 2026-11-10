package outputmodel

type UserOutput struct {
	ID       uint    `json:"-"`
	UUID     string  `json:"uuid"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	ImageURI *string `json:"imageUri,omitempty"`
}
