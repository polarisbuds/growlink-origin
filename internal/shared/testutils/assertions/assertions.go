package assertions

import (
	"errors"
	"strings"
	"testing"

	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func NoError(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err)
}

func ErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	require.ErrorIs(t, err, target)
}

func NotFound(t *testing.T, err error) {
	t.Helper()
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func NonZeroUUID(t *testing.T, id uuid.UUID) {
	t.Helper()
	assert.NotEqual(t, uuid.Nil, id)
}

func EqualUUID(t *testing.T, expected, actual uuid.UUID) {
	t.Helper()
	assert.Equal(t, expected, actual)
}

func OrderedIDs[T any](t *testing.T, items []T, getID func(T) uint) {
	t.Helper()
	for i := 1; i < len(items); i++ {
		prev := getID(items[i-1])
		curr := getID(items[i])
		assert.GreaterOrEqual(t, prev, curr, "expected id DESC ordering")
	}
}

func MeetingIDs(t *testing.T, expected, actual []entity.Meeting) {
	t.Helper()
	require.Len(t, actual, len(expected))
	for i := range expected {
		assert.Equal(t, expected[i].ID, actual[i].ID)
	}
}

func RequireMeetingPreloaded(t *testing.T, got *entity.Meeting, want *entity.Meeting) {
	t.Helper()
	require.NotNil(t, got)
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.UUID, got.UUID)
	assert.Equal(t, want.StatusID, got.StatusID)
	assert.Equal(t, want.TemplateID, got.TemplateID)
	assert.Equal(t, want.PeeringID, got.PeeringID)

	require.NotNil(t, got.Status)
	assert.Equal(t, uint(want.StatusID), got.Status.ID)

	require.NotNil(t, got.Template)
	assert.Equal(t, uint(want.TemplateID), got.Template.ID)

	require.NotNil(t, got.Peering)
	require.NotNil(t, got.Peering.Mentor)
	require.NotNil(t, got.Peering.Mentee)
}

func RequireMeetingNestedPreload(t *testing.T, got *entity.Meeting) {
	t.Helper()
	require.NotNil(t, got)
	require.NotNil(t, got.Template)
	require.NotEmpty(t, got.Template.Topics)
	require.NotEmpty(t, got.Answers)
	require.NotEmpty(t, got.Answers[0].AnswerComments)
	require.NotEmpty(t, got.MeetingComments)
	require.NotNil(t, got.MeetingComments[0].User)
}

func ConstraintViolation(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return
	}
	msg := strings.ToLower(err.Error())
	assert.True(t,
		strings.Contains(msg, "duplicate") ||
			strings.Contains(msg, "unique") ||
			strings.Contains(msg, "violates foreign key") ||
			strings.Contains(msg, "foreign key constraint"),
		"expected constraint violation, got: %v", err,
	)
}

func ForeignKeyViolation(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	msg := strings.ToLower(err.Error())
	assert.True(t,
		strings.Contains(msg, "violates foreign key") ||
			strings.Contains(msg, "foreign key constraint"),
		"expected foreign key violation, got: %v", err,
	)
}

func UniqueViolation(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return
	}
	msg := strings.ToLower(err.Error())
	assert.True(t,
		strings.Contains(msg, "duplicate") ||
			strings.Contains(msg, "unique"),
		"expected unique violation, got: %v", err,
	)
}
