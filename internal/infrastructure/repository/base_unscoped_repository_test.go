package repository_test

import (
	"context"
	"testing"

	"growlink/internal/domain/entity"
	"growlink/internal/infrastructure/repository"
	"growlink/internal/shared/testutils"
	"growlink/internal/shared/testutils/assertions"
	"growlink/internal/shared/testutils/fixtures"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newUnscopedMeetingRepo(t *testing.T, tx *gorm.DB) *repository.BaseUnScopedRepository[entity.Meeting] {
	t.Helper()
	return repository.NewBaseUnScopedRepository[entity.Meeting](context.Background(), tx)
}

func TestBaseUnScopedRepository_SoftDelete(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newBaseMeetingRepo(t, tx)
	unscoped := newUnscopedMeetingRepo(t, tx)
	meetingRepo := newMeetingRepo(t, tx)

	want := fixtures.Meeting(t, tx).Create()

	deleted, err := repo.Delete(1, "id = ?", want.ID)
	assertions.NoError(t, err)
	assert.Equal(t, 1, deleted)

	_, err = meetingRepo.GetByID(int64(want.ID))
	assertions.NotFound(t, err)

	got, err := unscoped.UFirst("id DESC", "id = ?", want.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, want.ID, got.ID)
}
