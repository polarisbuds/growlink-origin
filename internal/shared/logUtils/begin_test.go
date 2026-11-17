package logutils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

func captureLogger(buf *bytes.Buffer) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
}

func lastRecord(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) == 0 {
		t.Fatal("no log lines")
	}
	var rec map[string]any
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &rec); err != nil {
		t.Fatal(err)
	}
	return rec
}

func TestBegin_success(t *testing.T) {
	var buf bytes.Buffer
	captureLogger(&buf)

	ctx := WithRequestID(context.Background(), "req-1")
	done := Begin(ctx, LayerUsecase, "auth_usecase.LoginLDAP")
	done(nil)

	rec := lastRecord(t, &buf)
	if rec["msg"] != "処理終了" {
		t.Fatalf("msg: got %v", rec["msg"])
	}
	if rec["layer"] != LayerUsecase || rec["op"] != "auth_usecase.LoginLDAP" {
		t.Fatalf("layer/op: %+v", rec)
	}
	if rec["request_id"] != "req-1" {
		t.Fatalf("request_id: got %v", rec["request_id"])
	}
	if _, ok := rec["elapsed"]; !ok {
		t.Fatal("missing elapsed")
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	var start map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &start); err != nil {
		t.Fatal(err)
	}
	if start["msg"] != "処理開始" {
		t.Fatalf("start msg: got %v", start["msg"])
	}
}

func TestBegin_failure(t *testing.T) {
	var buf bytes.Buffer
	captureLogger(&buf)

	done := Begin(context.Background(), LayerRepository, "user_repository.GetByEmail")
	done(errors.New("db down"))

	rec := lastRecord(t, &buf)
	if rec["msg"] != "処理失敗" {
		t.Fatalf("msg: got %v", rec["msg"])
	}
	if rec["err"] != "db down" {
		t.Fatalf("err: got %v", rec["err"])
	}
}