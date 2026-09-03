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

func newPairingRepo(t *testing.T, tx *gorm.DB) *repository.PairingRepository {
	t.Helper()
	return repository.NewPairingRepository(context.Background(), tx)
}

func insertPairing(t *testing.T, tx *gorm.DB, mentor, mentee *entity.User) *entity.Pairing {
	t.Helper()
	pairing := &entity.Pairing{
		MentorID: mentor.ID,
		MenteeID: mentee.ID,
	}
	if err := tx.Create(pairing).Error; err != nil {
		t.Fatalf("failed to insert pairing: %v", err)
	}
	pairing.Mentor = mentor
	pairing.Mentee = mentee
	return pairing
}

func TestPairingRepository_CreatePairing(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newPairingRepo(t, tx)

	mentor := fixtures.User(t, tx, "create-mentor")
	mentee := fixtures.User(t, tx, "create-mentee")

	created, err := repo.CreatePairing(mentee.ID, mentor.ID)
	assertions.NoError(t, err)
	require.NotNil(t, created)
	assert.NotZero(t, created.ID)
	assert.Equal(t, mentor.ID, created.MentorID)
	assert.Equal(t, mentee.ID, created.MenteeID)

	got, err := repo.GetByMenteeID(mentee.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, mentor.ID, got.MentorID)
	assert.Equal(t, mentee.ID, got.MenteeID)
}

func TestPairingRepository_CreatePairing_ForeignKey(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newPairingRepo(t, tx)

	_, err := repo.CreatePairing(9_999_999, 9_999_998)
	assertions.ForeignKeyViolation(t, err)
}

func TestPairingRepository_GetByMenteeID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newPairingRepo(t, tx)

	mentor := fixtures.User(t, tx, "get-mentor")
	mentee := fixtures.User(t, tx, "get-mentee")
	want := insertPairing(t, tx, mentor, mentee)

	got, err := repo.GetByMenteeID(mentee.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.MentorID, got.MentorID)
	assert.Equal(t, want.MenteeID, got.MenteeID)
	assertions.NonZeroUUID(t, got.UUID)
}

func TestPairingRepository_GetByMenteeID_NotFound(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newPairingRepo(t, tx)

	got, err := repo.GetByMenteeID(9_999_999)
	assertions.NoError(t, err)
	assert.Nil(t, got)
}

func TestPairingRepository_GetByMentorID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newPairingRepo(t, tx)

	mentor := fixtures.User(t, tx, "list-mentor")
	mentee1 := fixtures.User(t, tx, "list-mentee1")
	mentee2 := fixtures.User(t, tx, "list-mentee2")
	first := insertPairing(t, tx, mentor, mentee1)
	second := insertPairing(t, tx, mentor, mentee2)

	got, err := repo.GetByMentorID(mentor.ID)
	assertions.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, first.ID, got[0].ID)
	assert.Equal(t, second.ID, got[1].ID)

	require.NotNil(t, got[0].Mentor)
	require.NotNil(t, got[0].Mentee)
	assert.Equal(t, mentor.ID, got[0].Mentor.ID)
	assert.Equal(t, mentee1.ID, got[0].Mentee.ID)

	require.NotNil(t, got[1].Mentor)
	require.NotNil(t, got[1].Mentee)
	assert.Equal(t, mentor.ID, got[1].Mentor.ID)
	assert.Equal(t, mentee2.ID, got[1].Mentee.ID)
}

func TestPairingRepository_BreakPairing(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newPairingRepo(t, tx)

	mentor := fixtures.User(t, tx, "break-mentor")
	mentee := fixtures.User(t, tx, "break-mentee")
	want := insertPairing(t, tx, mentor, mentee)

	err := repo.BreakPairing(want.ID)
	assertions.NoError(t, err)

	got, err := repo.GetByMenteeID(mentee.ID)
	assertions.NoError(t, err)
	assert.Nil(t, got)
}

func TestPairingRepository_BreakAllMenteePairing(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newPairingRepo(t, tx)

	mentorA := fixtures.User(t, tx, "breakall-mentor-a")
	mentorB := fixtures.User(t, tx, "breakall-mentor-b")
	mentee1 := fixtures.User(t, tx, "breakall-mentee1")
	mentee2 := fixtures.User(t, tx, "breakall-mentee2")
	mentee3 := fixtures.User(t, tx, "breakall-mentee3")
	insertPairing(t, tx, mentorA, mentee1)
	insertPairing(t, tx, mentorA, mentee2)
	kept := insertPairing(t, tx, mentorB, mentee3)

	err := repo.BreakAllMenteePairing(mentorA.ID)
	assertions.NoError(t, err)

	got1, err := repo.GetByMenteeID(mentee1.ID)
	assertions.NoError(t, err)
	assert.Nil(t, got1)

	got2, err := repo.GetByMenteeID(mentee2.ID)
	assertions.NoError(t, err)
	assert.Nil(t, got2)

	gotKept, err := repo.GetByMenteeID(mentee3.ID)
	assertions.NoError(t, err)
	require.NotNil(t, gotKept)
	assert.Equal(t, kept.ID, gotKept.ID)
}

func TestPairingRepository_GetAll(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newPairingRepo(t, tx)

	mentor := fixtures.User(t, tx, "getall-mentor")
	mentee1 := fixtures.User(t, tx, "getall-mentee1")
	mentee2 := fixtures.User(t, tx, "getall-mentee2")
	first := insertPairing(t, tx, mentor, mentee1)
	second := insertPairing(t, tx, mentor, mentee2)

	got, err := repo.GetAll()
	assertions.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, first.ID, got[0].ID)
	assert.Equal(t, second.ID, got[1].ID)

	for _, pairing := range got {
		require.NotNil(t, pairing.Mentor)
		require.NotNil(t, pairing.Mentee)
		assert.Equal(t, mentor.ID, pairing.Mentor.ID)
	}
}

func TestPairingRepository_GetAll_Empty(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newPairingRepo(t, tx)

	got, err := repo.GetAll()
	assertions.NoError(t, err)
	assert.Empty(t, got)
}
