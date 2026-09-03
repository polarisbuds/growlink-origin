package repository_test

import (
	"context"
	"testing"
	"time"

	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
	"growlink/internal/infrastructure/repository"
	"growlink/internal/shared/testutils"
	"growlink/internal/shared/testutils/assertions"
	"growlink/internal/shared/testutils/database"
	"growlink/internal/shared/testutils/fixtures"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func jst() *time.Location {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	return loc
}

func atJST(year int, month time.Month, day, hour, min, sec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, 0, jst())
}

func newMeetingRepo(t *testing.T, tx *gorm.DB) *repository.MeetingRepository {
	t.Helper()
	return repository.NewMeetingRepository(context.Background(), tx)
}

func TestMeetingRepository_InsertAndGetByID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newMeetingRepo(t, tx)

	want := fixtures.Meeting(t, tx).WithNestedPreloadGraph().Create()

	got, err := repo.GetByID(want.ID)
	assertions.NoError(t, err)
	assertions.RequireMeetingPreloaded(t, got, want)
	assertions.RequireMeetingNestedPreload(t, got)
}

func TestMeetingRepository_GetByUUID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newMeetingRepo(t, tx)

	want := fixtures.Meeting(t, tx).Create()
	assertions.NonZeroUUID(t, want.UUID)

	got, err := repo.GetByUUID(want.UUID)
	assertions.NoError(t, err)
	assertions.EqualUUID(t, want.UUID, got.UUID)
}

func TestMeetingRepository_GetAll_OrderingAndAssociation(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newMeetingRepo(t, tx)

	first := fixtures.Meeting(t, tx).Create()
	second := fixtures.Meeting(t, tx).Create()

	got, err := repo.GetAll()
	assertions.NoError(t, err)
	require.Len(t, got, 2)
	assertions.OrderedIDs(t, got, func(m entity.Meeting) uint { return m.ID })
	assertions.MeetingIDs(t, []entity.Meeting{*second, *first}, got)

	for _, meeting := range got {
		require.NotNil(t, meeting.Pairing)
		require.NotNil(t, meeting.Pairing.Mentor)
		require.NotNil(t, meeting.Pairing.Mentee)
		assert.True(t, meeting.Status.IsValid())
		require.NotNil(t, meeting.Template)
	}
}

func TestMeetingRepository_UniqueConstraint_UUID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newBaseMeetingRepo(t, tx)

	existing := fixtures.Meeting(t, tx).Create()

	err := repo.Insert(&entity.Meeting{
		UUID:          existing.UUID,
		ImplementedAt: time.Now(),
		Status:        existing.Status,
		TemplateID:    existing.TemplateID,
		PairingID:     existing.PairingID,
	})
	assertions.UniqueViolation(t, err)
}

func TestMeetingRepository_ForeignKeyConstraint(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newBaseMeetingRepo(t, tx)

	invalid := fixtures.MeetingWithoutFK(t, tx)
	err := repo.Insert(invalid)
	assertions.ForeignKeyViolation(t, err)
}

func TestMeetingRepository_UUIDAutoGeneration(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newMeetingRepo(t, tx)

	want := fixtures.Meeting(t, tx).Create()
	assertions.NonZeroUUID(t, want.UUID)

	got, err := repo.GetByUUID(want.UUID)
	assertions.NoError(t, err)
	assertions.EqualUUID(t, want.UUID, got.UUID)
}

func TestMeetingRepository_GetUntilToday_untilEndOfGivenDay(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newMeetingRepo(t, tx)

	base := atJST(2026, time.August, 15, 0, 0, 0)
	yesterday := atJST(2026, time.August, 14, 23, 59, 59)
	todayNoon := atJST(2026, time.August, 15, 12, 0, 0)
	tomorrow := atJST(2026, time.August, 16, 0, 0, 0)

	y := fixtures.Meeting(t, tx).WithImplementedAt(yesterday).Create()
	n := fixtures.Meeting(t, tx).WithImplementedAt(todayNoon).Create()
	_ = fixtures.Meeting(t, tx).WithImplementedAt(tomorrow).Create()

	got, err := repo.GetUntilToday(base)
	assertions.NoError(t, err)
	assertions.MeetingIDs(t, []entity.Meeting{*n, *y}, got)
}

func TestMeetingRepository_GetFromSpecificDateOnward_includesGivenDay(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newMeetingRepo(t, tx)

	base := atJST(2026, time.August, 15, 0, 0, 0)
	yesterday := atJST(2026, time.August, 14, 23, 59, 59)
	todayStart := atJST(2026, time.August, 15, 0, 0, 0)
	todayNoon := atJST(2026, time.August, 15, 12, 0, 0)
	tomorrow := atJST(2026, time.August, 16, 0, 0, 0)

	_ = fixtures.Meeting(t, tx).WithImplementedAt(yesterday).Create()
	start := fixtures.Meeting(t, tx).WithImplementedAt(todayStart).Create()
	noon := fixtures.Meeting(t, tx).WithImplementedAt(todayNoon).Create()
	next := fixtures.Meeting(t, tx).WithImplementedAt(tomorrow).Create()

	got, err := repo.GetFromSpecificDateOnward(base)
	assertions.NoError(t, err)
	assertions.MeetingIDs(t, []entity.Meeting{*next, *noon, *start}, got)
}

func TestMeetingRepository_GetBetweenDates_boundaries(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newMeetingRepo(t, tx)

	from := atJST(2026, time.August, 10, 15, 30, 0)
	to := atJST(2026, time.August, 12, 8, 0, 0)

	before := atJST(2026, time.August, 9, 23, 59, 59)
	start := atJST(2026, time.August, 10, 0, 0, 0)
	middle := atJST(2026, time.August, 11, 12, 0, 0)
	end := atJST(2026, time.August, 12, 23, 59, 59)
	after := atJST(2026, time.August, 13, 0, 0, 0)

	_ = fixtures.Meeting(t, tx).WithImplementedAt(before).Create()
	s := fixtures.Meeting(t, tx).WithImplementedAt(start).Create()
	m := fixtures.Meeting(t, tx).WithImplementedAt(middle).Create()
	e := fixtures.Meeting(t, tx).WithImplementedAt(end).Create()
	_ = fixtures.Meeting(t, tx).WithImplementedAt(after).Create()

	got, err := repo.GetBetweenDates(from, to)
	assertions.NoError(t, err)
	assertions.MeetingIDs(t, []entity.Meeting{*e, *m, *s}, got)
}

func TestMeetingRepository_GetBetweenDates_reversedArguments(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newMeetingRepo(t, tx)

	from := atJST(2026, time.August, 12, 0, 0, 0)
	to := atJST(2026, time.August, 10, 0, 0, 0)
	inRange := fixtures.Meeting(t, tx).WithImplementedAt(atJST(2026, time.August, 11, 12, 0, 0)).Create()

	got, err := repo.GetBetweenDates(from, to)
	assertions.NoError(t, err)
	assertions.MeetingIDs(t, []entity.Meeting{*inRange}, got)
}

func TestMeetingRepository_GetByPairingID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newMeetingRepo(t, tx)

	pairing := fixtures.Pairing(t, tx)
	otherPairing := fixtures.Pairing(t, tx)

	first := fixtures.Meeting(t, tx).WithPairing(pairing).Create()
	second := fixtures.Meeting(t, tx).WithPairing(pairing).Create()
	_ = fixtures.Meeting(t, tx).WithPairing(otherPairing).Create()

	got, err := repo.GetByPairingID(pairing.ID)
	assertions.NoError(t, err)
	require.Len(t, got, 2)
	assertions.MeetingIDs(t, []entity.Meeting{*first, *second}, got)

	for _, meeting := range got {
		require.NotNil(t, meeting.Pairing)
		require.NotNil(t, meeting.Pairing.Mentor)
		require.NotNil(t, meeting.Pairing.Mentee)
		assert.Nil(t, meeting.Template)
	}
}

func TestMeetingRepository_GetByPairingID_Empty(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newMeetingRepo(t, tx)

	pairing := fixtures.Pairing(t, tx)

	got, err := repo.GetByPairingID(pairing.ID)
	assertions.NoError(t, err)
	assert.Empty(t, got)
}

func TestMeetingRepository_Update(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newMeetingRepo(t, tx)

	want := fixtures.Meeting(t, tx).Create()
	want.Status = valueobject.MeetingStatusDone

	err := repo.Update(want.ID, want)
	assertions.NoError(t, err)

	got, err := repo.GetByID(want.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, valueobject.MeetingStatusDone, got.Status)
}

func TestMeetingRepository_UpdateImplementDate(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newMeetingRepo(t, tx)

	want := fixtures.Meeting(t, tx).Create()
	newDate := atJST(2026, time.September, 10, 15, 30, 0)

	err := repo.UpdateImplementDate(want.ID, newDate)
	assertions.NoError(t, err)

	got, err := repo.GetByID(want.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.ImplementedAt.Equal(newDate))
}

func TestMeetingRepository_TransactionRollbackIsolation(t *testing.T) {
	db, err := database.DB(context.Background())
	require.NoError(t, err)

	tx := db.Begin()
	require.NoError(t, tx.Error)

	repo := newMeetingRepo(t, tx)
	want := fixtures.Meeting(t, tx).Create()

	got, err := repo.GetByID(want.ID)
	assertions.NoError(t, err)
	require.NotNil(t, got)

	require.NoError(t, tx.Rollback().Error)

	otherTx := testutils.Begin(t)
	otherRepo := newMeetingRepo(t, otherTx)

	_, err = otherRepo.GetByID(want.ID)
	assertions.NotFound(t, err)
}
