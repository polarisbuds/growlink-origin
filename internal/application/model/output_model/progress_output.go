package outputmodel

type ProgressSummaryOutput struct {
	Pending int `json:"pending"`
	Draft   int `json:"draft"`
	Done    int `json:"done"`
	Total   int `json:"total"`
}
