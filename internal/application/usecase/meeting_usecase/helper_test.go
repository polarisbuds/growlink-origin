package meeting_usecase_test

import (
	"errors"
	"testing"

	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const (
	validMeetingID uint = 100
	validAnswerID  uint = 200
	validCommentID uint = 300
)

var (
	errGetAnswerFailed      = errors.New("get answer failed")
	errGetByMeetingFailed   = errors.New("get by meeting failed")
	errSubmitsFailed        = errors.New("submits failed")
	errUpdateStatusFailed   = errors.New("update status failed")
	errGetCommentFailed     = errors.New("get comment failed")
	errPostCommentFailed    = errors.New("post comment failed")
	errUpdateCommentFailed  = errors.New("update comment failed")
	errDeleteCommentFailed  = errors.New("delete comment failed")
)

func newTestAnswer(id, meetingID uint, content string) *entity.Answer {
	return &entity.Answer{
		Model:     gorm.Model{ID: id},
		UUID:      uuid.New(),
		Content:   content,
		MeetingID: meetingID,
		TopicID:   1,
		UserID:    1,
	}
}

func newTestAnswerComment(id, answerID uint, content string) *entity.AnswerComment {
	return &entity.AnswerComment{
		Model:    gorm.Model{ID: id},
		UUID:     uuid.New(),
		Content:  &content,
		AnswerID: answerID,
		UserID:   1,
	}
}

func newTestMeetingComment(id, meetingID uint, content string) *entity.MeetingComment {
	return &entity.MeetingComment{
		Model:     gorm.Model{ID: id},
		UUID:      uuid.New(),
		Content:   &content,
		MeetingID: meetingID,
		UserID:    1,
	}
}

func assertUsecaseErr(t *testing.T, gotErr error, wantErr bool, errContains, methodName string) {
	t.Helper()
	if gotErr != nil {
		if !wantErr {
			t.Errorf("%s() failed: %v", methodName, gotErr)
			return
		}
		if errContains != "" {
			assert.ErrorContains(t, gotErr, errContains)
		}
		return
	}
	if wantErr {
		t.Fatalf("%s() succeeded unexpectedly", methodName)
	}
}

func strPtr(s string) *string {
	return &s
}

func requireStringPtrEqual(t *testing.T, got, want *string) {
	t.Helper()
	require.NotNil(t, got)
	require.NotNil(t, want)
	assert.Equal(t, *want, *got)
}
