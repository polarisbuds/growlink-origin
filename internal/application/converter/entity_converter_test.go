package converter_test

import (
	"testing"
	"time"

	"growlink/internal/application/converter"
	inputmodel "growlink/internal/application/model/input_model"
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMeetingInputToEntityAndBack(t *testing.T) {
	implementedAt := time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC)
	input := inputmodel.MeetingInput{
		ImplementedAt: implementedAt,
		Status:        valueobject.MeetingStatusPending,
		TemplateID:    10,
		PairingID:     20,
	}

	got := converter.MeetingInputToEntity(input)
	require.NotNil(t, got)
	assert.True(t, got.ImplementedAt.Equal(implementedAt))
	assert.Equal(t, valueobject.MeetingStatusPending, got.Status)
	assert.Equal(t, uint(10), got.TemplateID)
	assert.Equal(t, uint(20), got.PairingID)

	entityMeeting := entity.Meeting{
		Model:         gorm.Model{ID: 1},
		ImplementedAt: implementedAt,
		Status:        valueobject.MeetingStatusPending,
		TemplateID:    10,
		PairingID:     20,
	}
	out := converter.MeetingEntityToOutput(entityMeeting)
	assert.Equal(t, uint(1), out.ID)
	assert.Equal(t, uint(10), out.Template.ID)
	assert.Equal(t, uint(20), out.Pairing.ID)
}

func TestTemplateInputToEntityAndBack(t *testing.T) {
	input := inputmodel.TemplateInput{
		Title: "weekly",
		Topics: []inputmodel.TopicInput{
			{DisplayNo: 1, Question: "Q1", IsRequire: true},
		},
	}

	got := converter.TemplateInputToEntity(input)
	require.NotNil(t, got)
	assert.Equal(t, "weekly", got.Title)
	require.Len(t, got.Topics, 1)
	assert.Equal(t, "Q1", got.Topics[0].Question)

	entityTemplate := entity.Template{
		Model: gorm.Model{ID: 5},
		Title: "weekly",
		Topics: []entity.Topic{
			{Model: gorm.Model{ID: 7}, DisplayNo: 1, Question: "Q1", IsRequire: true},
		},
	}
	out := converter.TemplateEntityToOutput(entityTemplate)
	assert.Equal(t, uint(5), out.ID)
	require.Len(t, out.Topics, 1)
	assert.Equal(t, uint(7), out.Topics[0].ID)
}

func TestUserInputToEntityAndBack(t *testing.T) {
	imageURI := "https://example.com/icon.png"
	input := inputmodel.UserInput{
		Name:         "Alice",
		Email:        "alice@example.com",
		ImageURI:     &imageURI,
		RoleID:       1,
		DepartmentID: 2,
		Authority:    valueobject.AuthorityOwner,
	}

	got := converter.UserInputToEntity(input)
	require.NotNil(t, got)
	assert.Equal(t, "Alice", got.Name)
	assert.Equal(t, "alice@example.com", got.Email)
	assert.Equal(t, &imageURI, got.ImageURI)

	entityUser := entity.User{
		Model: gorm.Model{ID: 99},
		Name:  "Alice",
		Email: "alice@example.com",
	}
	out := converter.UserEntityToOutput(entityUser)
	assert.Equal(t, uint(99), out.ID)
	assert.Equal(t, "Alice", out.Name)
}

func TestAnswerInputToEntityAndBackWithNilRelations(t *testing.T) {
	input := inputmodel.AnswerInput{
		Content:   "answer body",
		MeetingID: 3,
		TopicID:   4,
		UserID:    5,
	}

	got := converter.AnswerInputToEntity(input)
	require.NotNil(t, got)
	assert.Equal(t, "answer body", got.Content)
	assert.Equal(t, uint(3), got.MeetingID)

	entityAnswer := entity.Answer{
		Model:     gorm.Model{ID: 8},
		Content:   "answer body",
		MeetingID: 3,
		TopicID:   4,
		UserID:    5,
	}
	out := converter.AnswerEntityToOutput(entityAnswer)
	assert.Equal(t, uint(8), out.ID)
	assert.Equal(t, uint(3), out.Meeting.ID)
	assert.Equal(t, uint(4), out.Topic.ID)
}

func TestPairingEntityToOutputWithUsers(t *testing.T) {
	entityPairing := entity.Pairing{
		Model:    gorm.Model{ID: 11},
		MentorID: 21,
		MenteeID: 22,
		Mentor:   &entity.User{Model: gorm.Model{ID: 21}, Name: "mentor", Email: "mentor@example.com"},
		Mentee:   &entity.User{Model: gorm.Model{ID: 22}, Name: "mentee", Email: "mentee@example.com"},
	}

	out := converter.PairingEntityToOutput(entityPairing)
	assert.Equal(t, uint(11), out.ID)
	assert.Equal(t, "mentor", out.Mentor.Name)
	assert.Equal(t, "mentee", out.Mentee.Name)
}
