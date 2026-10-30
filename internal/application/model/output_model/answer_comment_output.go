package outputmodel

type AnswerCommentOutput struct {
	ID       uint          `json:"id"`
	Content  *string       `json:"content,omitempty"`
	AnswerID uint          `json:"answerId"`
	UserID   uint          `json:"userId"`
	Answer   AnswerOutput  `json:"answer"`
	User     UserOutput    `json:"user"`
}
