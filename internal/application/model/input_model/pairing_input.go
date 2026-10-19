package inputmodel

type PairingInput struct {
	MentorID      uint   `json:"mentorId"`
	MenteeID      uint   `json:"menteeId"`
	RecurWeek     *int   `json:"recurWeek"`
	RecurWeekday  *int   `json:"recurWeekday"`
	SharerUserIDs []uint `json:"sharerUserIds"`
}
