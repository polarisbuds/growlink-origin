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
	require.NotNil(t, got.ImplementedAt)
	assert.True(t, got.ImplementedAt.Equal(implementedAt))
	assert.Equal(t, valueobject.MeetingStatusPending, got.Status)
	assert.Equal(t, uint(10), got.TemplateID)
	assert.Equal(t, uint(20), got.PairingID)

	entityMeeting := entity.Meeting{
		Model:         gorm.Model{ID: 1},
		ImplementedAt: &implementedAt,
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
	iconFile := "icon.png"
	entityUser.ImageURI = &iconFile
	out := converter.UserEntityToOutput(entityUser)
	assert.Equal(t, uint(99), out.ID)
	assert.Equal(t, "Alice", out.Name)
	require.NotNil(t, out.ImageURI)
	assert.Equal(t, iconFile, *out.ImageURI)
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

func TestDepartmentHierarchyInputToEntityAndBack(t *testing.T) {
	input := inputmodel.DepartmentHierarchyInput{
		ParentDepartmentID: 1,
		ChildDepartmentID:  2,
	}

	got := converter.DepartmentHierarchyInputToEntity(input)
	require.NotNil(t, got)
	assert.Equal(t, uint(1), got.ParentDepartmentID)
	assert.Equal(t, uint(2), got.ChildDepartmentID)

	entityHierarchy := entity.DepartmentHierarchy{
		Model:              gorm.Model{ID: 3},
		ParentDepartmentID: 1,
		ChildDepartmentID:  2,
	}
	out := converter.DepartmentHierarchyEntityToOutput(entityHierarchy)
	assert.Equal(t, uint(3), out.ID)
	assert.Equal(t, uint(1), out.ParentDepartment.ID)
	assert.Equal(t, uint(2), out.ChildDepartment.ID)

	entityHierarchy.ParentDepartment = &entity.Department{Model: gorm.Model{ID: 1}, Name: "parent"}
	entityHierarchy.ChildDepartment = &entity.Department{Model: gorm.Model{ID: 2}, Name: "child"}
	outWithRelations := converter.DepartmentHierarchyEntityToOutput(entityHierarchy)
	assert.Equal(t, "parent", outWithRelations.ParentDepartment.Name)
	assert.Equal(t, "child", outWithRelations.ChildDepartment.Name)
}

func TestDepartmentRoleInputToEntityAndBack(t *testing.T) {
	input := inputmodel.DepartmentRoleInput{
		Name:         "Manager",
		UserID:       10,
		DepartmentID: 20,
	}

	got := converter.DepartmentRoleInputToEntity(input)
	require.NotNil(t, got)
	assert.Equal(t, "Manager", got.Name)
	assert.Equal(t, uint(10), got.UserID)
	assert.Equal(t, uint(20), got.DepartmentID)

	entityRole := entity.DepartmentRole{
		Model:        gorm.Model{ID: 5},
		Name:         "Manager",
		UserID:       10,
		DepartmentID: 20,
		User:         &entity.User{Model: gorm.Model{ID: 10}, Name: "Alice", Email: "alice@example.com"},
		Department:   &entity.Department{Model: gorm.Model{ID: 20}, Name: "Engineering"},
	}
	out := converter.DepartmentRoleEntityToOutput(entityRole)
	assert.Equal(t, uint(5), out.ID)
	assert.Equal(t, "Alice", out.User.Name)
	assert.Equal(t, "Engineering", out.Department.Name)
}

func TestUserReferenceUserInputToEntityAndBack(t *testing.T) {
	input := inputmodel.UserReferenceUserInput{
		UserID:          1,
		ReferableUserID: 2,
	}

	got := converter.UserReferenceUserInputToEntity(input)
	require.NotNil(t, got)
	assert.Equal(t, uint(1), got.UserID)
	assert.Equal(t, uint(2), got.ReferableUserID)

	entityRef := entity.UserReferenceUser{
		Model:           gorm.Model{ID: 7},
		UserID:          1,
		ReferableUserID: 2,
		User:            &entity.User{Model: gorm.Model{ID: 1}, Name: "owner", Email: "owner@example.com"},
		ReferableUser:   &entity.User{Model: gorm.Model{ID: 2}, Name: "target", Email: "target@example.com"},
	}
	out := converter.UserReferenceUserEntityToOutput(entityRef)
	assert.Equal(t, uint(7), out.ID)
	assert.Equal(t, "owner", out.User.Name)
	assert.Equal(t, "target", out.ReferableUser.Name)
}

func TestUserReferenceTeamInputToEntityAndBack(t *testing.T) {
	input := inputmodel.UserReferenceTeamInput{
		UserID:          1,
		ReferableTeamID: 3,
	}

	got := converter.UserReferenceTeamInputToEntity(input)
	require.NotNil(t, got)
	assert.Equal(t, uint(1), got.UserID)
	assert.Equal(t, uint(3), got.ReferableTeamID)

	entityRef := entity.UserReferenceTeam{
		Model:           gorm.Model{ID: 8},
		UserID:          1,
		ReferableTeamID: 3,
		User:            &entity.User{Model: gorm.Model{ID: 1}, Name: "owner", Email: "owner@example.com"},
		ReferableTeam:   &entity.Team{Model: gorm.Model{ID: 3}, Name: "Alpha", DepartmentID: 4},
	}
	out := converter.UserReferenceTeamEntityToOutput(entityRef)
	assert.Equal(t, uint(8), out.ID)
	assert.Equal(t, "owner", out.User.Name)
	assert.Equal(t, "Alpha", out.ReferableTeam.Name)
	assert.Empty(t, out.ReferableTeam.TeamMembers)
}

func TestUserReferenceDepartmentInputToEntityAndBack(t *testing.T) {
	input := inputmodel.UserReferenceDepartmentInput{
		UserID:                1,
		ReferableDepartmentID: 5,
	}

	got := converter.UserReferenceDepartmentInputToEntity(input)
	require.NotNil(t, got)
	assert.Equal(t, uint(1), got.UserID)
	assert.Equal(t, uint(5), got.ReferableDepartmentID)

	entityRef := entity.UserReferenceDepartment{
		Model:                 gorm.Model{ID: 9},
		UserID:                1,
		ReferableDepartmentID: 5,
		User:                  &entity.User{Model: gorm.Model{ID: 1}, Name: "owner", Email: "owner@example.com"},
		ReferableDepartment:   &entity.Department{Model: gorm.Model{ID: 5}, Name: "Sales"},
	}
	out := converter.UserReferenceDepartmentEntityToOutput(entityRef)
	assert.Equal(t, uint(9), out.ID)
	assert.Equal(t, "owner", out.User.Name)
	assert.Equal(t, "Sales", out.ReferableDepartment.Name)
}
