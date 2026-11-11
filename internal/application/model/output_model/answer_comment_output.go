package outputmodel

type AnswerCommentOutput struct {
	Content *string      `json:"content,omitempty"`
	Answer  AnswerOutput `json:"answer"`
	User    UserOutput   `json:"user"`
}
