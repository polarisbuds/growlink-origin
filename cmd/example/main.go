package main

import (
	"fmt"
	"growlink/internal/adapter/handler"
	"growlink/internal/infrastructure/external/ldap"
	infrastructure "growlink/internal/infrastructure/postgresql"
	"growlink/internal/shared/config"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	logDir                         = "logs"
	logFile, accessLogFile, awsLog *os.File
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	if _, err := os.Stat(logDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(logDir, 0755); err != nil {
				panic(err)
			}
		}
	}

	var err error
	logFile, err = os.OpenFile(time.Now().Format(logDir+`/log_20060102.log`), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}

	level := slog.LevelInfo
	if strings.ToLower(os.Getenv("MODE")[:2]) == "dev" {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	accessLogFile, err = os.OpenFile(time.Now().Format(logDir+`/log_20060102_access.log`), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		slog.Error("failed to load timezone", "err", err)
		os.Exit(1)
	}
	time.Local = loc
	maxProc := runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	slog.Info("Use processor limit", "maxProc", maxProc)
}

const (
	DATA_ROOT = "uploads"
)

func main() {
	defer func() {
		if err := accessLogFile.Close(); err != nil {
			slog.Error("failed to close access log file", "err", err)
			os.Exit(1)
		}

		if err := logFile.Close(); err != nil {
			panic(err.Error())
		}
	}()

	ldapSetting, err := ldap.NewLDAPSetting()
	if err != nil {
		slog.Error("LDAP設定情報の作成に失敗しました", slog.Any("msg", err))
		os.Exit(3)
	}
	ldapClient := ldap.NewClient(ldapSetting)
	user, err := ldapClient.Authenticate("ryota.kashio", "P@ssw0rd@@@@@@@@@@@")
	if err != nil {
		slog.Error("ADアカウント認証に失敗しました", slog.Any("msg", err))
		os.Exit(2)
	}
	slog.Info("ADアカウント認証に成功しました", slog.Any("user", *user))

	accessLog := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, accessLogFile), nil))
	panicLog := slog.New(slog.NewJSONHandler(accessLogFile, nil))
	_ = handler.NewEcho(accessLog, panicLog)

	dsn := config.GetPostgresqlDSN()
	fmt.Print("dsn: ", dsn)
	_, err = infrastructure.NewPostgresqlDB(dsn)
	if err != nil {
		slog.Error("failed to connect to postgresql", "err", err)
		os.Exit(1)
	}

	fmt.Println("SUCCESS.")
	os.Exit(0)
}
