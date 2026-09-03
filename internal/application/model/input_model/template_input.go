package inputmodel

type TemplateInput struct {
	Title  string       `json:"title"`
	Topics []TopicInput `json:"topics"`
}
