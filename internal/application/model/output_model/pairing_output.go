package outputmodel

type PairingOutput struct {
	ID            uint       `json:"id"`
	Mentor        UserOutput `json:"mentor"`
	Mentee        UserOutput `json:"mentee"`
	RecurWeek     *int       `json:"recurWeek"`
	RecurWeekday  *int       `json:"recurWeekday"`
	SharerUserIDs []uint     `json:"sharerUserIds"`
}
