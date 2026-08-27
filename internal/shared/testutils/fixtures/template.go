package fixtures

import (
	"fmt"
	"testing"
	"time"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// createTemplate inserts a Template with owner user (Layer2).
func createTemplate(t *testing.T, tx *gorm.DB) *entity.Template {
	t.Helper()

	user := createUser(t, tx, "template-owner")
	template := &entity.Template{
		Title:  fmt.Sprintf("template-%d", time.Now().UnixNano()),
		UserID: int64(user.ID),
	}
	if err := tx.Create(template).Error; err != nil {
		t.Fatalf("failed to create template fixture: %v", err)
	}
	template.User = user
	return template
}

// createTopic inserts a Topic for the given template (Layer3).
func createTopic(t *testing.T, tx *gorm.DB, template *entity.Template) *entity.Topic {
	t.Helper()

	topic := &entity.Topic{
		DisplayNo:  1,
		Question:   "test question",
		IsRequire:  true,
		TemplateID: int64(template.ID),
	}
	if err := tx.Create(topic).Error; err != nil {
		t.Fatalf("failed to create topic fixture: %v", err)
	}
	return topic
}
