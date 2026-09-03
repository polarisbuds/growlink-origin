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
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newTemplateRepo(t *testing.T, tx *gorm.DB) *repository.TemplateRepository {
	t.Helper()
	return repository.NewTemplateRepository(context.Background(), tx)
}

func TestTemplateRepository_GetAll(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTemplateRepo(t, tx)

	first := fixtures.Template(t, tx)
	fixtures.Topic(t, tx, first)
	second := fixtures.Template(t, tx)

	got, err := repo.GetAll()
	assertions.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, first.ID, got[0].ID)
	assert.Equal(t, second.ID, got[1].ID)

	require.NotNil(t, got[0].User)
	assert.Equal(t, first.UserID, got[0].User.ID)
	require.NotEmpty(t, got[0].Topics)
	assert.Equal(t, first.ID, got[0].Topics[0].TemplateID)

	require.NotNil(t, got[1].User)
	assert.Empty(t, got[1].Topics)
}

func TestTemplateRepository_Get(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTemplateRepo(t, tx)

	want := fixtures.Template(t, tx)
	topic := fixtures.Topic(t, tx, want)

	got, err := repo.Get(want.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Title, got.Title)
	assert.Equal(t, want.UserID, got.UserID)

	require.NotNil(t, got.User)
	assert.Equal(t, want.UserID, got.User.ID)
	require.Len(t, got.Topics, 1)
	assert.Equal(t, topic.ID, got.Topics[0].ID)
	assert.Equal(t, topic.Question, got.Topics[0].Question)
}

func TestTemplateRepository_Get_NotFound(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTemplateRepo(t, tx)

	_, err := repo.Get(9_999_999)
	assertions.NotFound(t, err)
}

func TestTemplateRepository_Create(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTemplateRepo(t, tx)

	owner := fixtures.User(t, tx, "template-create")
	title := fmt.Sprintf("created-template-%d", time.Now().UnixNano())
	input := &entity.Template{
		Title:  title,
		UserID: owner.ID,
	}

	err := repo.Create(input)
	assertions.NoError(t, err)
	assert.NotZero(t, input.ID)
	assertions.NonZeroUUID(t, input.UUID)

	got, err := repo.Get(input.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, title, got.Title)
	assert.Equal(t, owner.ID, got.UserID)
	require.NotNil(t, got.User)
	assert.Equal(t, owner.ID, got.User.ID)
}

func TestTemplateRepository_Create_ForeignKey(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTemplateRepo(t, tx)

	err := repo.Create(&entity.Template{
		Title:  fmt.Sprintf("invalid-template-%d", time.Now().UnixNano()),
		UserID: 9_999_999,
	})
	assertions.ForeignKeyViolation(t, err)
}

func TestTemplateRepository_Create_UniqueTitle(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTemplateRepo(t, tx)

	existing := fixtures.Template(t, tx)
	owner := fixtures.User(t, tx, "template-unique")

	err := repo.Create(&entity.Template{
		Title:  existing.Title,
		UserID: owner.ID,
	})
	assertions.UniqueViolation(t, err)
}

func TestTemplateRepository_Update(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTemplateRepo(t, tx)

	want := fixtures.Template(t, tx)
	updatedTitle := fmt.Sprintf("updated-template-%d", time.Now().UnixNano())
	want.Title = updatedTitle

	err := repo.Update(want.ID, want)
	assertions.NoError(t, err)

	got, err := repo.Get(want.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, updatedTitle, got.Title)
}

func TestTemplateRepository_Delete(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTemplateRepo(t, tx)

	want := fixtures.Template(t, tx)

	err := repo.Delete(want.ID)
	assertions.NoError(t, err)

	_, err = repo.Get(want.ID)
	assertions.NotFound(t, err)
}
