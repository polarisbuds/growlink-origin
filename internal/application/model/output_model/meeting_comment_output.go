package outputmodel

import "time"

type MeetingCommentOutput struct {
	ID        uint          `json:"id"`
	Content   string        `json:"content"`
	IsLimited bool          `json:"isLimited"`
	CreatedAt time.Time     `json:"createdAt"`
	MeetingID uint          `json:"meetingId"`
	UserID    uint          `json:"userId"`
	Meeting   MeetingOutput `json:"meeting"`
	User      UserOutput    `json:"user"`
}
