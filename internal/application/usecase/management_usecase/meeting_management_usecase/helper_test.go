package meeting_mng_usecase_test

import (
	"errors"
	"testing"
	"time"

	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"

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
	pairingID2        uint = 102
	meetingID1        uint = 201
	meetingID2        uint = 202
	meetingID3        uint = 203
	validTemplateID   uint = 301
)

var (
	errGetPairingsFailed   = errors.New("get pairings failed")
	errGetMeetingsFailed     = errors.New("get meetings failed")
	errGetRelationshipsFailed = errors.New("get relationships failed")
	errGetTeamMembersFailed  = errors.New("get team members failed")
	errGetUsersFailed        = errors.New("get users failed")
	errCreateMeetingFailed   = errors.New("create meeting failed")
	errUpdateMeetingFailed   = errors.New("update meeting failed")
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

func newTestMeeting(id, pairingID uint) entity.Meeting {
	return entity.Meeting{
		Model:         gorm.Model{ID: id},
		PairingID:     pairingID,
		TemplateID:    validTemplateID,
		Status:        valueobject.MeetingStatusPending,
		ImplementedAt: time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC),
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

func meetingOutputIDSet(meetings []outputmodel.MeetingOutput) map[uint]struct{} {
	set := make(map[uint]struct{}, len(meetings))
	for _, m := range meetings {
		set[m.ID] = struct{}{}
	}
	return set
}

func assertMeetingOutputIDSet(t *testing.T, wantIDs []uint, got []outputmodel.MeetingOutput) {
	t.Helper()
	gotSet := meetingOutputIDSet(got)
	assert.Len(t, got, len(wantIDs))
	for _, id := range wantIDs {
		_, ok := gotSet[id]
		assert.True(t, ok, "expected meeting ID %d in result", id)
	}
}
