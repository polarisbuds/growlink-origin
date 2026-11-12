package inputmodel

type CreateMeetingInput struct {
	MentorUuid      string   `json:"mentorUuid"`
	MentorID        uint     `json:"-"`
	MenteeUuids     []string `json:"menteeUuids"`
	MenteeIDs       []uint   `json:"-"`
	TemplateUuid    string   `json:"templateUuid"`
	TemplateID      uint     `json:"-"`
	ImplementedAt   string   `json:"implementedAt"`
	SharerUserUuids []string `json:"sharerUserUuids"`
	SharerUserIDs   []uint   `json:"-"`
}

type CreateMeetingByAllInput struct {
	TemplateUuid      string   `json:"templateUuid"`
	TemplateID        uint     `json:"-"`
	ImplementedAt     string   `json:"implementedAt"`
	SharerUserUuids   []string `json:"sharerUserUuids"`
	SharerUserIDs     []uint   `json:"-"`
	PairingUuids      []string `json:"pairingUuids"`
	PairingIDs        []uint   `json:"-"`
}

type CreateMeetingMonthlyInput struct {
	MentorUuid      string   `json:"mentorUuid"`
	MentorID        uint     `json:"-"`
	MenteeUuids     []string `json:"menteeUuids"`
	MenteeIDs       []uint   `json:"-"`
	TemplateUuid    string   `json:"templateUuid"`
	TemplateID      uint     `json:"-"`
	From            string   `json:"from"`
	To              string   `json:"to"`
	Week            int      `json:"week"`
	Weekday         int      `json:"weekday"`
	SharerUserUuids []string `json:"sharerUserUuids"`
	SharerUserIDs   []uint   `json:"-"`
}

type CreateMeetingMonthlyByPairingsInput struct {
	TemplateUuid      string   `json:"templateUuid"`
	TemplateID        uint     `json:"-"`
	From              string   `json:"from"`
	To                string   `json:"to"`
	Week              int      `json:"week"`
	Weekday           int      `json:"weekday"`
	SharerUserUuids   []string `json:"sharerUserUuids"`
	SharerUserIDs     []uint   `json:"-"`
	PairingUuids      []string `json:"pairingUuids"`
	PairingIDs        []uint   `json:"-"`
}
