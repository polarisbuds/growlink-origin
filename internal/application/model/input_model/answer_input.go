package inputmodel

type AnswerInput struct {
	Content   string `json:"content"`
	MeetingID uint   `json:"meetingId"`
	TopicID   uint   `json:"topicId"`
	UserID    uint   `json:"userId"`
}
