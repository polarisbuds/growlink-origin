package fixtures

import (
	"testing"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// MeetingComment inserts a MeetingComment with default dependencies.
func MeetingComment(t *testing.T, tx *gorm.DB, meeting *entity.Meeting, user *entity.User) *entity.MeetingComment {
	t.Helper()
	return createMeetingComment(t, tx, meeting, user)
}

func createMeetingComment(t *testing.T, tx *gorm.DB, meeting *entity.Meeting, user *entity.User) *entity.MeetingComment {
	t.Helper()

	if meeting == nil {
		meeting = Meeting(t, tx).Create()
	}
	if user == nil {
		user = createUser(t, tx, "meeting-comment")
	}

	meetingComment := &entity.MeetingComment{
		Content:   "meeting comment",
		MeetingID: meeting.ID,
		UserID:    user.ID,
	}
	if err := tx.Create(meetingComment).Error; err != nil {
		t.Fatalf("failed to create meeting comment fixture: %v", err)
	}
	meetingComment.Meeting = meeting
	meetingComment.User = user
	return meetingComment
}
