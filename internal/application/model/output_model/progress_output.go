package outputmodel

type ProgressSummaryOutput struct {
	Pending     int     `json:"pending"`
	Draft       int     `json:"draft"`
	Submitted   int     `json:"submitted"`
	Done        int     `json:"done"`
	Total       int     `json:"total"`
	DonePercent float64 `json:"donePercent"`
}

type ProgressBreakdownRowOutput struct {
	UUID           string   `json:"uuid"`
	Name           string   `json:"name"`
	DepartmentName string   `json:"departmentName,omitempty"`
	Pending        int      `json:"pending"`
	Draft          int      `json:"draft"`
	Submitted      int      `json:"submitted"`
	Done           int      `json:"done"`
	DonePercent    float64  `json:"donePercent"`
	UsersPending   []string `json:"usersPending"`
	UsersDraft     []string `json:"usersDraft"`
	UsersSubmitted []string `json:"usersSubmitted"`
	UsersDone      []string `json:"usersDone"`
}

type ProgressDoneMonthOutput struct {
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Label string `json:"label"`
	Done  int    `json:"done"`
}
