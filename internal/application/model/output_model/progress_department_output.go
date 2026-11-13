package outputmodel

type ProgressDepartmentTableOutput struct {
	Rows []ProgressDepartmentRowOutput `json:"rows"`
}

type ProgressDepartmentRowOutput struct {
	ProgressRowOutput
	Teams []ProgressTeamRowOutput `json:"teams"`
}

type ProgressRowOutput struct {
	UUID             string  `json:"uuid"`
	Name             string  `json:"name"`
	Depth            int     `json:"depth"`
	IsTotal          bool    `json:"isTotal"`
	UserCount        int     `json:"userCount"`
	PairCount        int     `json:"pairCount"`
	NotDoneCount     int     `json:"notDoneCount"`
	SubmittedCount   int     `json:"submittedCount"`
	SubmittedPercent float64 `json:"submittedPercent"`
	DoneCount        int     `json:"doneCount"`
	DonePercent      float64 `json:"donePercent"`
}
