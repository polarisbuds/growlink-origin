package outputmodel

type TemplateOutput struct {
	ID        uint          `json:"id"`
	UUID      string        `json:"uuid"`
	Title     string        `json:"title"`
	DisplayNo int           `json:"displayNo"`
	IsSample  bool          `json:"isSample"`
	Topics    []TopicOutput `json:"topics"`
}
