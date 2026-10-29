package repository_test

import (
	"context"
	"testing"

	"growlink/internal/domain/valueobject"
	"growlink/internal/infrastructure/repository"
	"growlink/internal/shared/testutils"
	"growlink/internal/shared/testutils/assertions"
	"growlink/internal/shared/testutils/fixtures"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestUserRelationshipRepository_GetByBossID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserRelationshipRepo(t, tx)

	subordinate := fixtures.User(t, tx, "getbyboss-sub")
	boss := fixtures.User(t, tx, "getbyboss-boss")
	assertions.NoError(t, repo.Create(subordinate.ID, boss.ID, valueobject.RelationshipTypeBoss))

	got, err := repo.GetByBossID(boss.ID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, subordinate.ID, got[0].SourceUserID)
	assert.Equal(t, boss.ID, got[0].TargetUserID)
}

func TestUserRelationshipRepository_DeleteByBossUserID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newUserRelationshipRepo(t, tx)

	sub1 := fixtures.User(t, tx, "delboss-sub1")
	sub2 := fixtures.User(t, tx, "delboss-sub2")
	boss := fixtures.User(t, tx, "delboss-boss")
	assertions.NoError(t, repo.Create(sub1.ID, boss.ID, valueobject.RelationshipTypeBoss))
	assertions.NoError(t, repo.Create(sub2.ID, boss.ID, valueobject.RelationshipTypeBoss))

	assertions.NoError(t, repo.DeleteByBossUserID(boss.ID))

	got, err := repo.GetByBossID(boss.ID)
	assertions.NoError(t, err)
	assert.Empty(t, got)
}
