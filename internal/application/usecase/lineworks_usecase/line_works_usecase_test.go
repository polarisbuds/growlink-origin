package lineworks_usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type fakeSender struct {
	ids  []string
	text string
	fail map[string]error
}

func (f *fakeSender) SendUserMessage(_ context.Context, userID, text string) error {
	f.ids = append(f.ids, userID)
	f.text = text
	if err, ok := f.fail[userID]; ok {
		return err
	}
	return nil
}

func TestSendMessageMultiple(t *testing.T) {
	sender := &fakeSender{}
	uc := New(sender, nil, nil, nil)
	if err := uc.SendMessage(context.Background(), []string{" a ", "", "b"}, "  hello  "); err != nil {
		t.Fatal(err)
	}
	if strings.Join(sender.ids, ",") != "a,b" || sender.text != "hello" {
		t.Fatalf("ids=%v text=%q", sender.ids, sender.text)
	}
}

func TestSendMessageContinuesAfterFailure(t *testing.T) {
	sender := &fakeSender{fail: map[string]error{"a": errors.New("down")}}
	uc := New(sender, nil, nil, nil)
	err := uc.SendMessage(context.Background(), []string{"a", "b"}, "hello")
	if err == nil || !strings.Contains(err.Error(), "a: down") {
		t.Fatal(err)
	}
	if strings.Join(sender.ids, ",") != "a,b" {
		t.Fatalf("ids=%v", sender.ids)
	}
}

func TestSendMessageRejectsEmpty(t *testing.T) {
	uc := New(&fakeSender{}, nil, nil, nil)
	if err := uc.SendMessage(context.Background(), []string{"a"}, "  "); err == nil {
		t.Fatal("expected empty message")
	}
	if err := uc.SendMessage(context.Background(), []string{"", " "}, "hello"); err == nil {
		t.Fatal("expected empty ids")
	}
}
