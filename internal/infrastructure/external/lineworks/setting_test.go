package lineworks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingRequired(t *testing.T) {
	t.Setenv(envClientID, "")
	t.Setenv(envClientSecret, "")
	t.Setenv(envServiceAccount, "")
	t.Setenv(envPrivateKeyPath, "")
	t.Setenv(envBotID, "")

	cfg, ok, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if ok || cfg.ClientID != "" {
		t.Fatalf("ok=%v cfg=%+v", ok, cfg)
	}
}

func TestLoadUnreadableKey(t *testing.T) {
	t.Setenv(envClientID, "cid")
	t.Setenv(envClientSecret, "secret")
	t.Setenv(envServiceAccount, "sa")
	t.Setenv(envBotID, "bot")
	t.Setenv(envPrivateKeyPath, filepath.Join(t.TempDir(), "missing.pem"))

	_, ok, err := Load()
	if err == nil || ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
}

func TestLoadDefaults(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "key.pem")
	if err := os.WriteFile(keyPath, []byte(`"-----BEGIN RSA PRIVATE KEY-----\nline1\n-----END RSA PRIVATE KEY-----\n"`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(envClientID, "cid")
	t.Setenv(envClientSecret, " secret ")
	t.Setenv(envServiceAccount, "sa")
	t.Setenv(envBotID, "bot")
	t.Setenv(envPrivateKeyPath, keyPath)
	t.Setenv(envAPIBase, "")
	t.Setenv(envTokenURL, "")
	t.Setenv(envScope, "")

	cfg, ok, err := Load()
	if err != nil || !ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	if cfg.APIBase != defaultAPIBase || cfg.TokenURL != defaultTokenURL || cfg.Scope != defaultScope {
		t.Fatalf("defaults: %+v", cfg)
	}
	if cfg.ClientSecret != " secret " {
		t.Fatalf("secret=%q", cfg.ClientSecret)
	}
	wantKey := "-----BEGIN RSA PRIVATE KEY-----\nline1\n-----END RSA PRIVATE KEY-----\n"
	if cfg.PrivateKey != wantKey {
		t.Fatalf("key=%q", cfg.PrivateKey)
	}
}
