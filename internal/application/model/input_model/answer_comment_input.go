package inputmodel

type AnswerCommentInput struct {
	Content  *string `json:"content,omitempty"`
	AnswerID uint    `json:"answerId"`
	UserID   uint    `json:"userId"`
}
