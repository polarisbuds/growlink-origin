package outputmodel

type AnswerOutput struct {
	ID             uint                `json:"id"`
	Content        string              `json:"content"`
	MeetingID      uint                `json:"meetingId"`
	TopicID        uint                `json:"topicId"`
	UserID         uint                `json:"userId"`
	Meeting        MeetingOutput       `json:"meeting"`
	Topic          TopicOutput         `json:"topic"`
	User           UserOutput          `json:"user"`
	AnswerComments []AnswerCommentOutput `json:"answerComments"`
}
