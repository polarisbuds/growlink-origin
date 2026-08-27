package fixtures

import (
	"fmt"
	"testing"
	"time"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// UserRelationship inserts a UserRelationship with default dependencies.
func UserRelationship(t *testing.T, tx *gorm.DB) *entity.UserRelationship {
	t.Helper()
	return createUserRelationship(t, tx)
}

// UserRelationshipType inserts a UserRelationshipType master record.
func UserRelationshipType(t *testing.T, tx *gorm.DB, prefix string) *entity.UserRelationshipType {
	t.Helper()
	return createUserRelationshipType(t, tx, prefix)
}

func createUserRelationshipType(t *testing.T, tx *gorm.DB, prefix string) *entity.UserRelationshipType {
	t.Helper()

	relationType := &entity.UserRelationshipType{
		Name: fmt.Sprintf("%s-relation-type-%d", prefix, time.Now().UnixNano()),
	}
	if err := tx.Create(relationType).Error; err != nil {
		t.Fatalf("failed to create user relationship type fixture: %v", err)
	}
	return relationType
}

func createUserRelationship(t *testing.T, tx *gorm.DB) *entity.UserRelationship {
	t.Helper()

	relationType := createUserRelationshipType(t, tx, "relationship")
	sourceUser := createUser(t, tx, "source")
	targetUser := createUser(t, tx, "target")

	relationship := &entity.UserRelationship{
		RelationTypeID: int64(relationType.ID),
		SourceUserID:   int64(sourceUser.ID),
		TargetUserID:   int64(targetUser.ID),
	}
	if err := tx.Create(relationship).Error; err != nil {
		t.Fatalf("failed to create user relationship fixture: %v", err)
	}
	relationship.RelationType = relationType
	relationship.SourceUser = sourceUser
	relationship.TargetUser = targetUser
	return relationship
}
