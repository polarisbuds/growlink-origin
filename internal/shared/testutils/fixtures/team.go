package fixtures

import (
	"fmt"
	"testing"
	"time"

	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

// Team inserts a Team with default dependencies.
func Team(t *testing.T, tx *gorm.DB) *entity.Team {
	t.Helper()
	return createTeam(t, tx)
}

// TeamMember inserts a TeamMember with default dependencies.
func TeamMember(t *testing.T, tx *gorm.DB, team *entity.Team, member *entity.User) *entity.TeamMember {
	t.Helper()
	return createTeamMember(t, tx, team, member)
}

func createTeam(t *testing.T, tx *gorm.DB) *entity.Team {
	t.Helper()

	department := createDepartment(t, tx, "team")
	team := &entity.Team{
		Name:         fmt.Sprintf("team-%d", time.Now().UnixNano()),
		DepartmentID: int64(department.ID),
	}
	if err := tx.Create(team).Error; err != nil {
		t.Fatalf("failed to create team fixture: %v", err)
	}
	team.Department = department
	return team
}

func createTeamMember(t *testing.T, tx *gorm.DB, team *entity.Team, member *entity.User) *entity.TeamMember {
	t.Helper()

	if team == nil {
		team = createTeam(t, tx)
	}
	if member == nil {
		member = createUser(t, tx, "team-member")
	}

	teamMember := &entity.TeamMember{
		TeamID:   int64(team.ID),
		MemberID: int64(member.ID),
	}
	if err := tx.Create(teamMember).Error; err != nil {
		t.Fatalf("failed to create team member fixture: %v", err)
	}
	teamMember.Team = team
	teamMember.Member = member
	return teamMember
}
