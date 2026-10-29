package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"growlink/internal/domain/entity"
	"growlink/internal/infrastructure/repository"
	"growlink/internal/shared/testutils"
	"growlink/internal/shared/testutils/assertions"
	"growlink/internal/shared/testutils/fixtures"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func newTopicRepo(t *testing.T, tx *gorm.DB) *repository.TopicRepository {
	t.Helper()
	return repository.NewTopicRepository(context.Background(), tx)
}

func TestTopicRepository_Create(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTopicRepo(t, tx)

	template := fixtures.Template(t, tx)
	input := &entity.Topic{
		DisplayNo:  2,
		Question:   fmt.Sprintf("created-topic-%d", time.Now().UnixNano()),
		IsRequire:  false,
		TemplateID: template.ID,
	}

	err := repo.Create(input)
	assertions.NoError(t, err)
	assert.NotZero(t, input.ID)
	assertions.NonZeroUUID(t, input.UUID)

	var got entity.Topic
	err = tx.First(&got, input.ID).Error
	assertions.NoError(t, err)
	assert.Equal(t, input.DisplayNo, got.DisplayNo)
	assert.Equal(t, input.Question, got.Question)
	assert.Equal(t, input.IsRequire, got.IsRequire)
	assert.Equal(t, template.ID, got.TemplateID)
}

func TestTopicRepository_Create_ForeignKey(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTopicRepo(t, tx)

	err := repo.Create(&entity.Topic{
		DisplayNo:  1,
		Question:   "invalid-topic",
		IsRequire:  true,
		TemplateID: 9_999_999,
	})
	assertions.ForeignKeyViolation(t, err)
}
