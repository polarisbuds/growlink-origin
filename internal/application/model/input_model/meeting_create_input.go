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
	MentorID      uint   `json:"mentorId"`
	MenteeIDs     []uint `json:"menteeIds"`
	TemplateID    uint   `json:"templateId"`
	From          string `json:"from"`
	To            string `json:"to"`
	Week          int    `json:"week"`
	Weekday       int    `json:"weekday"`
	SharerUserIDs []uint `json:"sharerUserIds"`
}
