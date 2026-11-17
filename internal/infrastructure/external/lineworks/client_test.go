package lineworks

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func testKey(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}))
}

func TestSendUserMessage(t *testing.T) {
	const secret = "super-secret-value"
	var tokenCalls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/token") {
			tokenCalls.Add(1)
			if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
				t.Errorf("content-type=%s", r.Header.Get("Content-Type"))
			}
			body, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(body), "client_secret="+secret) {
				t.Errorf("token form missing secret")
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"tok","expires_in":3600}`))
			return
		}
		if r.URL.Path != "/v1.0/bots/bot/users/user@example/messages" {
			t.Errorf("path=%s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer tok" {
			t.Errorf("auth=%s", r.Header.Get("Authorization"))
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"content":{"text":"hello","type":"text"}}` {
			t.Errorf("body=%s", body)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	client := testClient(t, secret, srv.URL)
	if err := client.SendUserMessage(context.Background(), "user@example", "hello"); err != nil {
		t.Fatal(err)
	}
	if err := client.SendUserMessage(context.Background(), "user@example", "hello"); err != nil {
		t.Fatal(err)
	}
	if tokenCalls.Load() != 1 {
		t.Fatalf("token calls=%d", tokenCalls.Load())
	}
}

func TestSendUserMessageRefreshesOnceOnUnauthorized(t *testing.T) {
	var tokenCalls atomic.Int32
	var sends atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/token") {
			n := tokenCalls.Add(1)
			_, _ = w.Write([]byte(`{"access_token":"tok` + string(rune('0'+n)) + `","expires_in":"3600"}`))
			return
		}
		if sends.Add(1) == 1 {
			http.Error(w, "super-secret-value", http.StatusUnauthorized)
			return
		}
		if r.Header.Get("Authorization") != "Bearer tok2" {
			t.Errorf("auth=%s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := testClient(t, "super-secret-value", srv.URL)
	if err := client.SendUserMessage(context.Background(), "tanaka", "hello"); err != nil {
		t.Fatal(err)
	}
	if tokenCalls.Load() != 2 || sends.Load() != 2 {
		t.Fatalf("token=%d sends=%d", tokenCalls.Load(), sends.Load())
	}
}

func TestSendUserMessageErrorOmitsSecretAndBody(t *testing.T) {
	const secret = "super-secret-value"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/token") {
			_, _ = w.Write([]byte(`{"access_token":"tok","expires_in":3600}`))
			return
		}
		http.Error(w, secret, http.StatusNotFound)
	}))
	defer srv.Close()

	client := testClient(t, secret, srv.URL)
	err := client.SendUserMessage(context.Background(), "tanaka", "hello")
	if err == nil || !strings.Contains(err.Error(), "status=404") || strings.Contains(err.Error(), secret) {
		t.Fatal(err)
	}
}

func testClient(t *testing.T, secret, base string) *Client {
	t.Helper()
	client, err := NewClient(Setting{
		ClientID:       "cid",
		ClientSecret:   secret,
		ServiceAccount: "sa",
		PrivateKey:     testKey(t),
		BotID:          "bot",
		APIBase:        base,
		TokenURL:       base + "/token",
		Scope:          "bot.message",
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}
