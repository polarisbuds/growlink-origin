package fixtures

import (
	"testing"

	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"

	"gorm.io/gorm"
)

// UserRelationship inserts a UserRelationship with default dependencies.
func UserRelationship(t *testing.T, tx *gorm.DB) *entity.UserRelationship {
	t.Helper()
	return createUserRelationship(t, tx)
}

func createUserRelationship(t *testing.T, tx *gorm.DB) *entity.UserRelationship {
	t.Helper()

	sourceUser := createUser(t, tx, "source")
	targetUser := createUser(t, tx, "target")

	relationship := &entity.UserRelationship{
		RelationType: valueobject.RelationshipTypeMentor,
		SourceUserID: sourceUser.ID,
		TargetUserID: targetUser.ID,
	}
	if err := tx.Create(relationship).Error; err != nil {
		t.Fatalf("failed to create user relationship fixture: %v", err)
	}
	relationship.SourceUser = sourceUser
	relationship.TargetUser = targetUser
	return relationship
}
