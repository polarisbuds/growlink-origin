package repository_test

import (
	"context"
	"testing"

	"growlink/internal/infrastructure/repository"
	"growlink/internal/shared/testutils"
	"growlink/internal/shared/testutils/assertions"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func newUserRelationshipRepo(t *testing.T, tx *gorm.DB) *repository.UserRelationshipRepository {
	t.Helper()
	return repository.NewUserRelationshipRepository(context.Background(), tx)
}

func TestUserRelationshipRepository_Get_NotFound(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserRelationshipRepo(t, tx)

	got, err := repo.Get(9_999_999)
	assertions.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, got)
}
