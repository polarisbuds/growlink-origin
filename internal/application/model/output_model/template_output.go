package outputmodel

type TemplateOutput struct {
	ID     uint          `json:"id"`
	Title  string        `json:"title"`
	Topics []TopicOutput `json:"topics"`
}
