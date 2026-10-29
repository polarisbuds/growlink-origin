package fixtures

import (
	"testing"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// Answer inserts an Answer with default dependencies.
func Answer(t *testing.T, tx *gorm.DB, meeting *entity.Meeting, topic *entity.Topic, user *entity.User) *entity.Answer {
	t.Helper()
	return createAnswer(t, tx, meeting, topic, user)
}

// AnswerComment inserts an AnswerComment with default dependencies.
func AnswerComment(t *testing.T, tx *gorm.DB, answer *entity.Answer, user *entity.User) *entity.AnswerComment {
	t.Helper()
	return createAnswerComment(t, tx, answer, user)
}

func createAnswer(t *testing.T, tx *gorm.DB, meeting *entity.Meeting, topic *entity.Topic, user *entity.User) *entity.Answer {
	t.Helper()

	if meeting == nil {
		meeting = Meeting(t, tx).Create()
	}
	if topic == nil {
		template := createTemplate(t, tx)
		topic = createTopic(t, tx, template)
	}
	if user == nil {
		user = createUser(t, tx, "answer")
	}

	answer := &entity.Answer{
		Content:   "test answer",
		MeetingID: meeting.ID,
		TopicID:   topic.ID,
		UserID:    user.ID,
	}
	if err := tx.Create(answer).Error; err != nil {
		t.Fatalf("failed to create answer fixture: %v", err)
	}
	answer.Meeting = meeting
	answer.Topic = topic
	answer.User = user
	return answer
}

func createAnswerComment(t *testing.T, tx *gorm.DB, answer *entity.Answer, user *entity.User) *entity.AnswerComment {
	t.Helper()

	if answer == nil {
		answer = createAnswer(t, tx, nil, nil, nil)
	}
	if user == nil {
		user = createUser(t, tx, "answer-comment")
	}

	commentContent := "answer comment"
	answerComment := &entity.AnswerComment{
		Content:  &commentContent,
		AnswerID: answer.ID,
		UserID:   user.ID,
	}
	if err := tx.Create(answerComment).Error; err != nil {
		t.Fatalf("failed to create answer comment fixture: %v", err)
	}
	answerComment.Answer = answer
	answerComment.User = user
	return answerComment
}
