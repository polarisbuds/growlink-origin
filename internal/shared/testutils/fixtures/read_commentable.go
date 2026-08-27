package fixtures

import (
	"testing"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// ReadCommentable inserts a ReadCommentable with default dependencies.
func ReadCommentable(t *testing.T, tx *gorm.DB, meeting *entity.Meeting, user *entity.User) *entity.ReadCommentable {
	t.Helper()
	return createReadCommentable(t, tx, meeting, user)
}

func createReadCommentable(t *testing.T, tx *gorm.DB, meeting *entity.Meeting, user *entity.User) *entity.ReadCommentable {
	t.Helper()

	if meeting == nil {
		meeting = Meeting(t, tx).Create()
	}
	if user == nil {
		user = createUser(t, tx, "read-commentable")
	}

	readCommentable := &entity.ReadCommentable{
		MeetingID: int64(meeting.ID),
		UserID:    int64(user.ID),
	}
	if err := tx.Create(readCommentable).Error; err != nil {
		t.Fatalf("failed to create read commentable fixture: %v", err)
	}
	readCommentable.Meeting = meeting
	readCommentable.User = user
	return readCommentable
}
