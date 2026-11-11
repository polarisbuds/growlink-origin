package outputmodel

type TemplateOutput struct {
	ID          uint          `json:"-"`
	UUID        string        `json:"uuid"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	IsSample    bool          `json:"isSample"`
	IsLock      bool          `json:"isLock"`
	IsDefault   bool          `json:"isDefault"`
	CreatorName string        `json:"creatorName"`
	CreatedAt   string        `json:"createdAt"`
	Topics      []TopicOutput `json:"topics"`
}
