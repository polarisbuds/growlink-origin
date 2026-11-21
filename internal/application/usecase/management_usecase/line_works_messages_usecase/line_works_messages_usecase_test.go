package line_works_messages_usecase

import (
	"testing"

	"github.com/google/uuid"
)

func validMessages() []LineWorksMessageUpdateItem {
	return []LineWorksMessageUpdateItem{
		{UUID: uuid.NewString(), Content: "[year]年[month]月のログが未提出です。"},
		{UUID: uuid.NewString(), Content: "[mentee]さんが[month]月のログを提出しました。"},
		{UUID: uuid.NewString(), Content: "[year]年[month]月の面談が未実施です。"},
	}
}

func TestValidateUpdateInput(t *testing.T) {
	if err := validateUpdateInput(UpdateLineWorksMessagesInput{Messages: validMessages()}); err != nil {
		t.Fatalf("expected ok: %v", err)
	}
	if err := validateUpdateInput(UpdateLineWorksMessagesInput{Messages: validMessages()[:2]}); err == nil {
		t.Fatal("expected count error")
	}
	msgs := validMessages()
	msgs[0].Content = "   "
	if err := validateUpdateInput(UpdateLineWorksMessagesInput{Messages: msgs}); err == nil {
		t.Fatal("expected empty content error")
	}
	msgs = validMessages()
	msgs[1].UUID = "not-a-uuid"
	if err := validateUpdateInput(UpdateLineWorksMessagesInput{Messages: msgs}); err == nil {
		t.Fatal("expected invalid uuid error")
	}
	msgs = validMessages()
	msgs[2].UUID = msgs[0].UUID
	if err := validateUpdateInput(UpdateLineWorksMessagesInput{Messages: msgs}); err == nil {
		t.Fatal("expected duplicate uuid error")
	}
}

func TestLineWorksMessages_Update_validationOnly(t *testing.T) {
	uc := &LineWorksMessages{}
	if err := uc.Update(1, UpdateLineWorksMessagesInput{Messages: []LineWorksMessageUpdateItem{
		{UUID: uuid.NewString(), Content: "x"},
	}}); err == nil {
		t.Fatal("expected count error before repo")
	}
}