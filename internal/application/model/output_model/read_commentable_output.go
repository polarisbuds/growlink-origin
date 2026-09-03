package outputmodel

type ReadCommentableOutput struct {
	ID        uint          `json:"id"`
	MeetingID uint          `json:"meetingId"`
	UserID    uint          `json:"userId"`
	Meeting   MeetingOutput `json:"meeting"`
	User      UserOutput    `json:"user"`
}
