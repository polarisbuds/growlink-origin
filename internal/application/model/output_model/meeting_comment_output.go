package outputmodel

import "time"

type MeetingCommentOutput struct {
	UUID      string        `json:"uuid"`
	Content   string        `json:"content"`
	IsLimited bool          `json:"isLimited"`
	CreatedAt time.Time     `json:"createdAt"`
	Meeting   MeetingOutput `json:"meeting"`
	User      UserOutput    `json:"user"`
}
