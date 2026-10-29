package inputmodel

type MeetingCommentInput struct {
	Content   string `json:"content"`
	MeetingID uint   `json:"meetingId"`
	UserID    uint   `json:"userId"`
}
