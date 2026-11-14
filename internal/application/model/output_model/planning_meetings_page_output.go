package outputmodel

type PlanningMeetingsPageOutput struct {
	Items    []MeetingOutput `json:"items"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
}