package template_mng_usecase_test

import (
	"errors"
	"testing"

	"growlink/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	validTemplateID uint = 301
	topicID1        uint = 401
	topicID2        uint = 402
)

var (
	errCreateTopicFailed    = errors.New("create topic failed")
	errCreateTemplateFailed = errors.New("create template failed")
	errGetTemplatesFailed   = errors.New("get templates failed")
	errGetTemplateFailed    = errors.New("get template failed")
	errUpdateTemplateFailed = errors.New("update template failed")
	errDeleteTemplateFailed = errors.New("delete template failed")
)

func newTestTopic(id uint, displayNo uint, question string, isRequire bool) entity.Topic {
	return entity.Topic{
		Model:     gorm.Model{ID: id},
		DisplayNo: displayNo,
		Question:  question,
		IsRequire: isRequire,
	}
}

func newTestTemplate(id uint, title string, topics ...entity.Topic) entity.Template {
	return entity.Template{
		Model:  gorm.Model{ID: id},
		Title:  title,
		Topics: topics,
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
