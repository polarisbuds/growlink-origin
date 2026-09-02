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

	_, err = meetingRepo.GetByID(want.ID)
	assertions.NotFound(t, err)

	got, err := unscoped.UFirst("id DESC", "id = ?", want.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, want.ID, got.ID)
}

func TestBaseUnScopedRepository_UChain(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newBaseMeetingRepo(t, tx)
	unscoped := newUnscopedMeetingRepo(t, tx)
	ctx := context.Background()

	first := fixtures.Meeting(t, tx).Create()
	second := fixtures.Meeting(t, tx).Create()

	deleted, err := repo.Delete(1, "id = ?", first.ID)
	assertions.NoError(t, err)
	assert.Equal(t, 1, deleted)

	got, err := unscoped.UChain("id DESC").Find(ctx)
	assertions.NoError(t, err)
	assertions.MeetingIDs(t, []entity.Meeting{*second, *first}, got)
	for _, meeting := range got {
		assert.Nil(t, meeting.Template)
		assert.Nil(t, meeting.Pairing)
	}

	preloaded, err := unscoped.UChain("id DESC", "Template", "Pairing", "Pairing.Mentor", "Pairing.Mentee").
		Where("id = ?", first.ID).First(ctx)
	assertions.NoError(t, err)
	assertions.RequireMeetingPreloaded(t, &preloaded, first)
}
