package outputmodel

type PairingOutput struct {
	ID              uint         `json:"id"`
	UUID            string       `json:"uuid"`
	Mentor          UserOutput   `json:"mentor"`
	Mentee          UserOutput   `json:"mentee"`
	Mentees         []UserOutput `json:"mentees"`
	RecurWeek       *int         `json:"recurWeek"`
	RecurWeekday    *int         `json:"recurWeekday"`
	SharerUserIDs   []uint       `json:"sharerUserIds"`
	CreatedByUserID uint         `json:"createdByUserId"`
}
