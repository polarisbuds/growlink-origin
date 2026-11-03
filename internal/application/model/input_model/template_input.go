package inputmodel

type TemplateInput struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Topics      []TopicInput `json:"topics"`
}
