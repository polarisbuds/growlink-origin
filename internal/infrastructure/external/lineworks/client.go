package lineworks

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenRefreshSkew = 60 * time.Second

// Client は LINE WORKS Bot API へユーザー宛テキストを送る。
type Client struct {
	httpClient *http.Client
	cfg        Setting
	key        *rsa.PrivateKey

	mu     sync.Mutex
	token  string
	expiry time.Time
}

func NewClient(cfg Setting) (*Client, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("LINE WORKSの秘密鍵を読み取れません")
	}
	if cfg.APIBase == "" {
		cfg.APIBase = defaultAPIBase
	}
	if cfg.TokenURL == "" {
		cfg.TokenURL = defaultTokenURL
	}
	if cfg.Scope == "" {
		cfg.Scope = defaultScope
	}
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cfg:        cfg,
		key:        key,
	}, nil
}

// SendUserMessage は userID のユーザーへテキストを1通送る。
func (c *Client) SendUserMessage(ctx context.Context, userID, text string) error {
	endpoint := fmt.Sprintf("%s/v1.0/bots/%s/users/%s/messages",
		strings.TrimRight(c.cfg.APIBase, "/"),
		url.PathEscape(c.cfg.BotID),
		url.PathEscape(userID),
	)
	status, err := c.postMessage(ctx, endpoint, text, true)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("line works send failed: status=%d", status)
	}
	return nil
}

func (c *Client) postMessage(ctx context.Context, endpoint, text string, allowRefresh bool) (int, error) {
	token, err := c.accessToken(ctx)
	if err != nil {
		return 0, fmt.Errorf("line works token: %w", err)
	}
	payload, err := json.Marshal(map[string]any{
		"content": map[string]string{
			"type": "text",
			"text": text,
		},
	})
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(payload)))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode == http.StatusUnauthorized && allowRefresh {
		c.clearToken()
		return c.postMessage(ctx, endpoint, text, false)
	}
	return resp.StatusCode, nil
}

func (c *Client) accessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && time.Now().Before(c.expiry) {
		return c.token, nil
	}
	token, expiry, err := c.fetchToken(ctx)
	if err != nil {
		return "", err
	}
	c.token = token
	c.expiry = expiry
	return token, nil
}

func (c *Client) clearToken() {
	c.mu.Lock()
	c.token = ""
	c.expiry = time.Time{}
	c.mu.Unlock()
}

func (c *Client) fetchToken(ctx context.Context) (string, time.Time, error) {
	assertion, err := c.signAssertion()
	if err != nil {
		return "", time.Time{}, err
	}
	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Set("assertion", assertion)
	form.Set("client_id", c.cfg.ClientID)
	form.Set("client_secret", c.cfg.ClientSecret)
	form.Set("scope", c.cfg.Scope)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", time.Time{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", time.Time{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", time.Time{}, fmt.Errorf("line works token request failed: status=%d", resp.StatusCode)
	}
	var parsed struct {
		AccessToken string          `json:"access_token"`
		ExpiresIn   json.RawMessage `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil || parsed.AccessToken == "" {
		return "", time.Time{}, fmt.Errorf("line works token request failed: status=%d", resp.StatusCode)
	}
	expiresIn, err := parseExpiresIn(parsed.ExpiresIn)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid expires_in")
	}
	if expiresIn <= 0 {
		expiresIn = 3600
	}
	skew := tokenRefreshSkew
	if time.Duration(expiresIn)*time.Second <= skew*2 {
		skew = time.Duration(expiresIn) * time.Second / 2
	}
	expiry := time.Now().Add(time.Duration(expiresIn)*time.Second - skew)
	return parsed.AccessToken, expiry, nil
}

func parseExpiresIn(raw json.RawMessage) (int64, error) {
	raw = []byte(strings.TrimSpace(string(raw)))
	if len(raw) == 0 || string(raw) == "null" {
		return 3600, nil
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return 0, err
		}
		return strconv.ParseInt(s, 10, 64)
	}
	var n int64
	if err := json.Unmarshal(raw, &n); err != nil {
		return 0, err
	}
	return n, nil
}

func (c *Client) signAssertion() (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": c.cfg.ClientID,
		"sub": c.cfg.ServiceAccount,
		"iat": now.Unix(),
		"exp": now.Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(c.key)
}
