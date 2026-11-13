package outputmodel

type ProgressTeamTableOutput struct {
	Rows []ProgressDepartmentRowOutput `json:"rows"`
}

type ProgressTeamRowOutput struct {
	ProgressRowOutput
	DepartmentUUIDs []string `json:"departments"`
}
