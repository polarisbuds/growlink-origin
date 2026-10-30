package inputmodel

type CreateMeetingInput struct {
	MentorID      uint   `json:"mentorId"`
	MenteeIDs     []uint `json:"menteeIds"`
	TemplateID    uint   `json:"templateId"`
	ImplementedAt string `json:"implementedAt"`
	SharerUserIDs []uint `json:"sharerUserIds"`
}

type CreateMeetingByAllInput struct {
	TemplateID      uint   `json:"templateId"`
	ImplementedAt   string `json:"implementedAt"`
	SharerUserIDs   []uint `json:"sharerUserIds"`
	PairingIDs      []uint `json:"pairingIds"`
}

type CreateMeetingMonthlyInput struct {
	TemplateID    uint   `json:"templateId"`
	BaseDate      string `json:"baseDate"`
	Week          int    `json:"week"`
	Weekday       int    `json:"weekday"`
	SharerUserIDs []uint `json:"sharerUserIds"`
	PairingIDs    []uint `json:"pairingIds"`
}
