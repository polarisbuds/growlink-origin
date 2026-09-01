package setting_usecase_test

import (
	"errors"
	"strings"
	"testing"

	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const validUserID uint = 10

var (
	errUserNotFound      = errors.New("user not found")
	errNewFileFailed     = errors.New("new file failed")
	errUpdateFailed      = errors.New("update failed")
	errExistsCheckFailed = errors.New("exists check failed")
	errDeleteFailed      = errors.New("delete failed")
)

func newTestUser(id uint, imageURI *string) *entity.User {
	return &entity.User{
		Model:    gorm.Model{ID: uint(id)},
		Name:     "icon-user",
		Email:    "icon@example.com",
		ImageURI: imageURI,
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

func assertUUIDFilename(t *testing.T, filename, ext string) {
	t.Helper()
	require.True(t, strings.HasSuffix(filename, ext), "filename %q should end with %q", filename, ext)
	id, err := uuid.Parse(strings.TrimSuffix(filename, ext))
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}
