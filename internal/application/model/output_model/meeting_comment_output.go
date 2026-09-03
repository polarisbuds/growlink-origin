package outputmodel

type MeetingCommentOutput struct {
	ID        uint          `json:"id"`
	Content   *string       `json:"content,omitempty"`
	MeetingID uint          `json:"meetingId"`
	UserID    uint          `json:"userId"`
	Meeting   MeetingOutput `json:"meeting"`
	User      UserOutput    `json:"user"`
}
