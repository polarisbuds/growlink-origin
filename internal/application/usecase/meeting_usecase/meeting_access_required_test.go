package meeting_usecase

import (
	"errors"
	"testing"

	"growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRequiredAnswers(t *testing.T) {
	requiredUUID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	optionalUUID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	meeting := &entity.Meeting{
		Template: &entity.Template{
			Topics: []entity.Topic{
				{UUID: requiredUUID, IsRequire: true},
				{UUID: optionalUUID, IsRequire: false},
			},
		},
	}

	t.Run("必須がすべて埋まっている", func(t *testing.T) {
		err := validateRequiredAnswers(meeting, []outputmodel.AnswerOutput{
			{Content: "ok", Topic: outputmodel.TopicOutput{UUID: requiredUUID.String()}},
			{Content: "", Topic: outputmodel.TopicOutput{UUID: optionalUUID.String()}},
		})
		require.NoError(t, err)
	})

	t.Run("必須が空", func(t *testing.T) {
		err := validateRequiredAnswers(meeting, []outputmodel.AnswerOutput{
			{Content: "", Topic: outputmodel.TopicOutput{UUID: requiredUUID.String()}},
		})
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrRequireAnswerNotAnswered))
	})

	t.Run("必須の回答行がない", func(t *testing.T) {
		err := validateRequiredAnswers(meeting, []outputmodel.AnswerOutput{
			{Content: "optional", Topic: outputmodel.TopicOutput{UUID: optionalUUID.String()}},
		})
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrRequireAnswerNotAnswered))
	})

	t.Run("任意のみ空", func(t *testing.T) {
		err := validateRequiredAnswers(meeting, []outputmodel.AnswerOutput{
			{Content: "required", Topic: outputmodel.TopicOutput{UUID: requiredUUID.String()}},
		})
		require.NoError(t, err)
	})
}