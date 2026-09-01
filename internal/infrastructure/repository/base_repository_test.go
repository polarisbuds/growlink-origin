package repository_test

import (
	"context"
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

func newBaseMeetingRepo(t *testing.T, tx *gorm.DB) *repository.BaseRepository[entity.Meeting] {
	t.Helper()
	return repository.NewBaseRepository[entity.Meeting](context.Background(), tx)
}

func TestBaseRepository_CRUD(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newBaseMeetingRepo(t, tx)
	meetingRepo := newMeetingRepo(t, tx)

	base := fixtures.Meeting(t, tx).Create()

	toInsert := &entity.Meeting{
		ImplementedAt: base.ImplementedAt,
		StatusID:      base.StatusID,
		TemplateID:    base.TemplateID,
		PairingID:     base.PairingID,
	}
	err := repo.Insert(toInsert)
	assertions.NoError(t, err)

	got, err := repo.First("id DESC", "id = ?", toInsert.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, toInsert.ID, got.ID)

	updatedAt := atJST(2026, time.August, 1, 12, 0, 0)
	_, err = repo.Update("implemented_at", updatedAt, "id = ?", toInsert.ID)
	assertions.NoError(t, err)

	got, err = repo.First("id DESC", "id = ?", toInsert.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, updatedAt.Equal(got.ImplementedAt))

	deleted, err := repo.Delete(1, "id = ?", toInsert.ID)
	assertions.NoError(t, err)
	assert.Equal(t, 1, deleted)

	_, err = meetingRepo.GetByID(toInsert.ID)
	assertions.NotFound(t, err)
}
