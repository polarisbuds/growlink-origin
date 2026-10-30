package pairing_mng_usecase_test

import (
	"errors"
	"testing"

	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	validMentorID     uint = 21
	validBossID       uint = 31
	validTeamID       uint = 41
	validDepartmentID uint = 51
	validMenteeID     uint = 11
	validMemberID     uint = 12
	pairingID1        uint = 101
)

var (
	errGetPairingsFailed      = errors.New("get pairings failed")
	errGetRelationshipsFailed = errors.New("get relationships failed")
	errGetTeamMembersFailed   = errors.New("get team members failed")
	errGetUsersFailed         = errors.New("get users failed")
)

func newTestUser(id uint, name string) *entity.User {
	return &entity.User{
		Model: gorm.Model{ID: id},
		Name:  name,
		Email: name + "@example.com",
	}
}

func newTestPairing(id, menteeID, mentorID uint) *entity.Pairing {
	return &entity.Pairing{
		Model:    gorm.Model{ID: id},
		MenteeID: menteeID,
		MentorID: mentorID,
		Mentee:   newTestUser(menteeID, "mentee"),
		Mentor:   newTestUser(mentorID, "mentor"),
	}
}

func newTestUserRelationship(id, sourceUserID, targetUserID uint) entity.UserRelationship {
	return entity.UserRelationship{
		Model:        gorm.Model{ID: id},
		SourceUserID: sourceUserID,
		TargetUserID: targetUserID,
	}
}

func newTestTeamMember(id, teamID, memberID uint) entity.TeamMember {
	return entity.TeamMember{
		Model:    gorm.Model{ID: id},
		TeamID:   teamID,
		MemberID: memberID,
	}
}

func assertUsecaseErr(t *testing.T, gotErr error, wantErr bool, errContains, methodName string) {
	t.Helper()
	if gotErr != nil {
		if !wantErr {
			t.Errorf("%s() failed: %v", methodName, gotErr)
			return
		}
		if errContains != "" {
			assert.ErrorContains(t, gotErr, errContains)
		}
		return
	}
	if wantErr {
		t.Fatalf("%s() succeeded unexpectedly", methodName)
	}
}

func assertPairingOutputIDSet(t *testing.T, wantIDs []uint, got []outputmodel.PairingOutput) {
	t.Helper()
	gotSet := make(map[uint]struct{}, len(got))
	for _, p := range got {
		gotSet[p.ID] = struct{}{}
	}
	assert.Len(t, got, len(wantIDs))
	for _, id := range wantIDs {
		_, ok := gotSet[id]
		assert.True(t, ok, "expected pairing ID %d in result", id)
	}
}
