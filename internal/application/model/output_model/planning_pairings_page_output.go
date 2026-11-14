package outputmodel

type PlanningPairingsPageOutput struct {
	Items    []PairingOutput `json:"items"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
}