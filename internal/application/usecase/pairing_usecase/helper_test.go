package pairing_usecase_test

import (
	"errors"
	"testing"

	"growlink/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	validMenteeID uint = 11
	validMentorID uint = 21
	oldPairingID  uint = 101
)

var (
	errGetPairingFailed    = errors.New("get pairing failed")
	errCreatePairingFailed = errors.New("create pairing failed")
	errBreakPairingFailed  = errors.New("break pairing failed")
)

func newTestUser(id uint, name string) *entity.User {
	return &entity.User{
		Model: gorm.Model{ID: id},
		Name:  name,
		Email: name + "@example.com",
	}
}

func newTestPairing(id, menteeID, mentorID uint) *entity.Pairing {
	mentee := newTestUser(menteeID, "mentee")
	mentor := newTestUser(mentorID, "mentor")
	return &entity.Pairing{
		Model:    gorm.Model{ID: id},
		MenteeID: menteeID,
		MentorID: mentorID,
		Mentee:   mentee,
		Mentor:   mentor,
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
