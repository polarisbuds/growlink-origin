package outputmodel

type ProgressSummaryOutput struct {
	Pending     int     `json:"pending"`
	Draft       int     `json:"draft"`
	Done        int     `json:"done"`
	Total       int     `json:"total"`
	DonePercent float64 `json:"donePercent"`
}

type ProgressBreakdownRowOutput struct {
	ID             uint     `json:"id"`
	Name           string   `json:"name"`
	DepartmentName string   `json:"departmentName,omitempty"`
	Pending        int      `json:"pending"`
	Draft          int      `json:"draft"`
	Done           int      `json:"done"`
	DonePercent    float64  `json:"donePercent"`
	UsersPending   []string `json:"usersPending"`
	UsersDraft     []string `json:"usersDraft"`
	UsersDone      []string `json:"usersDone"`
}
