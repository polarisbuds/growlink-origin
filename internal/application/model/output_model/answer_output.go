package outputmodel

type AnswerOutput struct {
	ID             uint                  `json:"-"`
	UUID           string                `json:"uuid"`
	Content        string                `json:"content"`
	MeetingID      uint                  `json:"-"`
	TopicID        uint                  `json:"-"`
	UserID         uint                  `json:"-"`
	Meeting        MeetingOutput         `json:"meeting"`
	Topic          TopicOutput           `json:"topic"`
	User           UserOutput            `json:"user"`
	AnswerComments []AnswerCommentOutput `json:"answerComments"`
}
