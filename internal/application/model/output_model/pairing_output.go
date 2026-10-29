package outputmodel

type PairingTargetOutput struct {
	Kind     string `json:"kind"`
	TargetID uint   `json:"targetId"`
	Label    string `json:"label"`
}

type PairingOutput struct {
	ID              uint                  `json:"id"`
	UUID            string                `json:"uuid"`
	Mentor          UserOutput            `json:"mentor"`
	Mentee          UserOutput            `json:"mentee"`
	Targets         []PairingTargetOutput `json:"targets"`
	CreatedByUserID uint                  `json:"createdByUserId"`
}
