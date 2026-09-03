package inputmodel

type MeetingCommentInput struct {
	Content   *string `json:"content,omitempty"`
	MeetingID uint    `json:"meetingId"`
	UserID    uint    `json:"userId"`
}
