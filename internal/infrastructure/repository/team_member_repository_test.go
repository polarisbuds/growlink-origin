package repository_test

import (
	"context"
	"testing"

	"growlink/internal/infrastructure/repository"
	"growlink/internal/shared/testutils"
	"growlink/internal/shared/testutils/assertions"
	"growlink/internal/shared/testutils/fixtures"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newTeamMemberRepo(t *testing.T, tx *gorm.DB) *repository.TeamMemberRepository {
	t.Helper()
	return repository.NewTeamMemberRepository(context.Background(), tx)
}

func TestTeamMemberRepository_Create(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTeamMemberRepo(t, tx)

	team := fixtures.Team(t, tx)
	member := fixtures.User(t, tx, "create-member")

	err := repo.Create(team.ID, member.ID, true)
	assertions.NoError(t, err)

	got, err := repo.GetByTeamID(team.ID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, team.ID, got[0].TeamID)
	assert.Equal(t, member.ID, got[0].MemberID)
	assert.True(t, got[0].IsReader)
}

func TestTeamMemberRepository_Create_ForeignKey(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTeamMemberRepo(t, tx)

	err := repo.Create(9_999_999, 9_999_998, false)
	assertions.ForeignKeyViolation(t, err)
}

func TestTeamMemberRepository_Update(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTeamMemberRepo(t, tx)

	team := fixtures.Team(t, tx)
	member := fixtures.User(t, tx, "update-member")
	want := fixtures.TeamMemberWithReader(t, tx, team, member, false)

	want.IsReader = true
	err := repo.Update(want.ID, want)
	assertions.NoError(t, err)

	got, err := repo.GetByMemberID(member.ID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.True(t, got[0].IsReader)
}

func TestTeamMemberRepository_Delete(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTeamMemberRepo(t, tx)

	team := fixtures.Team(t, tx)
	member := fixtures.User(t, tx, "delete-member")
	want := fixtures.TeamMember(t, tx, team, member)

	err := repo.Delete(want.ID)
	assertions.NoError(t, err)

	got, err := repo.GetByTeamID(team.ID)
	assertions.NoError(t, err)
	assert.Empty(t, got)
}

func TestTeamMemberRepository_GetAll(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTeamMemberRepo(t, tx)

	first := fixtures.TeamMember(t, tx, nil, nil)
	second := fixtures.TeamMember(t, tx, nil, nil)

	got, err := repo.GetAll()
	assertions.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, first.ID, got[0].ID)
	assert.Equal(t, second.ID, got[1].ID)

	for _, tm := range got {
		require.NotNil(t, tm.Team)
		require.NotNil(t, tm.Team.Department)
		require.NotNil(t, tm.Member)
	}
}

func TestTeamMemberRepository_GetByTeamID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTeamMemberRepo(t, tx)

	team := fixtures.Team(t, tx)
	otherTeam := fixtures.Team(t, tx)
	want := fixtures.TeamMember(t, tx, team, nil)
	_ = fixtures.TeamMember(t, tx, otherTeam, nil)

	got, err := repo.GetByTeamID(team.ID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, want.ID, got[0].ID)
	require.NotNil(t, got[0].Team)
	assert.Equal(t, team.ID, got[0].Team.ID)
}

func TestTeamMemberRepository_GetByMemberID(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTeamMemberRepo(t, tx)

	member := fixtures.User(t, tx, "by-member")
	want := fixtures.TeamMember(t, tx, nil, member)
	_ = fixtures.TeamMember(t, tx, nil, nil)

	got, err := repo.GetByMemberID(member.ID)
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, want.ID, got[0].ID)
	require.NotNil(t, got[0].Member)
	assert.Equal(t, member.ID, got[0].Member.ID)
}

func TestTeamMemberRepository_GetReaders(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTeamMemberRepo(t, tx)

	team := fixtures.Team(t, tx)
	reader := fixtures.TeamMemberWithReader(t, tx, team, nil, true)
	_ = fixtures.TeamMemberWithReader(t, tx, team, nil, false)

	got, err := repo.GetReaders()
	assertions.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, reader.ID, got[0].ID)
	assert.True(t, got[0].IsReader)
}

func TestTeamMemberRepository_GetByTeamID_Empty(t *testing.T) {
	tx := testutils.Begin(t)
	repo := newTeamMemberRepo(t, tx)

	got, err := repo.GetByTeamID(9_999_999)
	assertions.NoError(t, err)
	assert.Empty(t, got)
}
