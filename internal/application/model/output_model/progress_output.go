package outputmodel

type ProgressSummaryOutput struct {
	Pending int `json:"pending"`
	Draft   int `json:"draft"`
	Done    int `json:"done"`
	Total   int `json:"total"`
}

type ProgressBreakdownRowOutput struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Pending int    `json:"pending"`
	Draft   int    `json:"draft"`
	Done    int    `json:"done"`
}
