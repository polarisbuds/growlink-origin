package inputmodel

type PairingInput struct {
	MentorID      uint   `json:"mentorId"`
	MenteeID      uint   `json:"menteeId"`
	MenteeIDs     []uint `json:"menteeIds"`
	RecurWeek     *int   `json:"recurWeek"`
	RecurWeekday  *int   `json:"recurWeekday"`
	SharerUserIDs []uint `json:"sharerUserIds"`
}

func (in PairingInput) ResolvedMenteeIDs() []uint {
	if len(in.MenteeIDs) > 0 {
		return in.MenteeIDs
	}
	if in.MenteeID != 0 {
		return []uint{in.MenteeID}
	}
	return nil
}
