package inputmodel

type CreateMeetingInput struct {
	PairingID       uint   `json:"pairingId"`
	TemplateID      uint   `json:"templateId"`
	ImplementedAt   string `json:"implementedAt"`
	SharerUserIDs   []uint `json:"sharerUserIds"`
}

type CreateMeetingByAllInput struct {
	TemplateID      uint   `json:"templateId"`
	ImplementedAt   string `json:"implementedAt"`
	SharerUserIDs   []uint `json:"sharerUserIds"`
}

type CreateMeetingMonthlyInput struct {
	TemplateID    uint   `json:"templateId"`
	BaseDate      string `json:"baseDate"`
	Week          int    `json:"week"`
	Weekday       int    `json:"weekday"`
	SharerUserIDs []uint `json:"sharerUserIds"`
}
