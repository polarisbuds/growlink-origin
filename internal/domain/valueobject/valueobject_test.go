package valueobject_test

import (
	"testing"

	"growlink/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthority_String_Parse_Int_UInt(t *testing.T) {
	tests := []struct {
		name   string
		value  valueobject.Authority
		str    string
		intVal int
	}{
		{name: "general", value: valueobject.AuthorityGeneral, str: "general", intVal: 1},
		{name: "boss", value: valueobject.AuthorityBoss, str: "boss", intVal: 2},
		{name: "owner", value: valueobject.AuthorityOwner, str: "owner", intVal: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.value.IsValid())
			assert.Equal(t, tt.str, tt.value.String())
			assert.Equal(t, tt.intVal, tt.value.Int())
			assert.Equal(t, uint(tt.intVal), tt.value.UInt())

			parsed, err := valueobject.ParseAuthority(tt.str)
			require.NoError(t, err)
			assert.Equal(t, tt.value, parsed)
		})
	}
}

func TestAuthority_Invalid(t *testing.T) {
	invalid := valueobject.Authority(0)
	assert.False(t, invalid.IsValid())
	assert.Equal(t, "", invalid.String())
	assert.Equal(t, 0, invalid.Int())
	assert.Equal(t, uint(0), invalid.UInt())

	_, err := valueobject.ParseAuthority("invalid")
	require.Error(t, err)
}

func TestMeetingStatus_String_Parse_Int_UInt(t *testing.T) {
	tests := []struct {
		name   string
		value  valueobject.MeetingStatus
		str    string
		intVal int
	}{
		{name: "pending", value: valueobject.MeetingStatusPending, str: "pending", intVal: 1},
		{name: "draft", value: valueobject.MeetingStatusDraft, str: "draft", intVal: 2},
		{name: "done", value: valueobject.MeetingStatusDone, str: "done", intVal: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.value.IsValid())
			assert.Equal(t, tt.str, tt.value.String())
			assert.Equal(t, tt.intVal, tt.value.Int())
			assert.Equal(t, uint(tt.intVal), tt.value.UInt())

			parsed, err := valueobject.ParseMeetingStatus(tt.str)
			require.NoError(t, err)
			assert.Equal(t, tt.value, parsed)
		})
	}
}

func TestMeetingStatus_Invalid(t *testing.T) {
	invalid := valueobject.MeetingStatus(0)
	assert.False(t, invalid.IsValid())
	assert.Equal(t, "", invalid.String())
	assert.Equal(t, 0, invalid.Int())
	assert.Equal(t, uint(0), invalid.UInt())

	_, err := valueobject.ParseMeetingStatus("invalid")
	require.Error(t, err)
}

func TestRelationshipType_String_Parse_Int_UInt(t *testing.T) {
	tests := []struct {
		name   string
		value  valueobject.RelationshipType
		str    string
		intVal int
	}{
		{name: "boss", value: valueobject.RelationshipTypeBoss, str: "boss", intVal: 1},
		{name: "subordinate", value: valueobject.RelationshipTypeSubordinate, str: "subordinate", intVal: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.value.IsValid())
			assert.Equal(t, tt.str, tt.value.String())
			assert.Equal(t, tt.intVal, tt.value.Int())
			assert.Equal(t, uint(tt.intVal), tt.value.UInt())

			parsed, err := valueobject.ParseRelationshipType(tt.str)
			require.NoError(t, err)
			assert.Equal(t, tt.value, parsed)
		})
	}
}

func TestRelationshipType_Invalid(t *testing.T) {
	invalid := valueobject.RelationshipType(0)
	assert.False(t, invalid.IsValid())
	assert.Equal(t, "", invalid.String())
	assert.Equal(t, 0, invalid.Int())
	assert.Equal(t, uint(0), invalid.UInt())

	_, err := valueobject.ParseRelationshipType("invalid")
	require.Error(t, err)
}
