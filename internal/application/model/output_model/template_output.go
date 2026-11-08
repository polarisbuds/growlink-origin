package outputmodel

type TemplateOutput struct {
	ID          uint          `json:"id"`
	UUID        string        `json:"uuid"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	DisplayNo   int           `json:"displayNo"`
	IsSample    bool          `json:"isSample"`
	IsLock      bool          `json:"isLock"`
	CreatorID   uint          `json:"creatorId"`
	CreatorName string        `json:"creatorName"`
	CreatedAt   string        `json:"createdAt"`
	Topics      []TopicOutput `json:"topics"`
}
