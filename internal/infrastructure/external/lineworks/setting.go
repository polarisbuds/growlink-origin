package lineworks

import (
	"fmt"
	"os"
	"strings"
)

const (
	envClientID       = "LINEWORKS_CLIENT_ID"
	envClientSecret   = "LINEWORKS_CLIENT_SECRET"
	envServiceAccount = "LINEWORKS_SERVICE_ACCOUNT"
	envPrivateKeyPath = "LINEWORKS_PRIVATE_KEY_PATH"
	envBotID          = "LINEWORKS_BOT_ID"
	envAPIBase        = "LINEWORKS_API_BASE"
	envTokenURL       = "LINEWORKS_TOKEN_URL"
	envScope          = "LINEWORKS_SCOPE"

	defaultAPIBase  = "https://www.worksapis.com"
	defaultTokenURL = "https://auth.worksmobile.com/oauth2/v2.0/token"
	defaultScope    = "bot.message"
)

// Setting は LINE WORKS Bot API へユーザー宛メッセージを送るための接続情報。
type Setting struct {
	ClientID       string
	ClientSecret   string
	ServiceAccount string
	PrivateKey     string
	BotID          string
	APIBase        string
	TokenURL       string
	Scope          string
}

// Load は送信に必要な環境変数が揃っているときだけ設定を返す。
// 必須項目が欠けているときは ok=false。秘密鍵ファイルが読めないときだけ error を返す。
func Load() (Setting, bool, error) {
	clientID := strings.TrimSpace(os.Getenv(envClientID))
	clientSecret := os.Getenv(envClientSecret)
	serviceAccount := strings.TrimSpace(os.Getenv(envServiceAccount))
	botID := strings.TrimSpace(os.Getenv(envBotID))
	keyPath := strings.TrimSpace(os.Getenv(envPrivateKeyPath))

	var privateKey string
	if keyPath != "" {
		privateKeyBytes, err := os.ReadFile(keyPath)
		if err != nil {
			return Setting{}, false, err
		}
		privateKey = normalizePEM(string(privateKeyBytes))
	}

	if clientID == "" || strings.TrimSpace(clientSecret) == "" || serviceAccount == "" || privateKey == "" || botID == "" {
		return Setting{}, false, fmt.Errorf("環境変数が正しく指定されていません: %s", fmt.Sprintf(`clientID="%s", strings.TrimSpace(clientSecret)="%s", serviceAccount="%s", privateKey="%s", botID="%s"`, clientID, strings.TrimSpace(clientSecret), serviceAccount, privateKey, botID))
	}

	apiBase := strings.TrimRight(strings.TrimSpace(os.Getenv(envAPIBase)), "/")
	if apiBase == "" {
		apiBase = defaultAPIBase
	}
	tokenURL := strings.TrimSpace(os.Getenv(envTokenURL))
	if tokenURL == "" {
		tokenURL = defaultTokenURL
	}
	scope := strings.TrimSpace(os.Getenv(envScope))
	if scope == "" {
		scope = defaultScope
	}

	return Setting{
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		ServiceAccount: serviceAccount,
		PrivateKey:     privateKey,
		BotID:          botID,
		APIBase:        apiBase,
		TokenURL:       tokenURL,
		Scope:          scope,
	}, true, nil
}

func normalizePEM(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.Trim(raw, `"`)
	return strings.ReplaceAll(raw, `\n`, "\n")
}
