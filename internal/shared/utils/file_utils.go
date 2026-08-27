package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	osExecutable = os.Executable
	osGetwd      = os.Getwd
)

func GetAbsoluteTargetDir(folderName string) (string, error) {
	// 実行ファイルの絶対パスを取得
	exePath, err := osExecutable()
	if err != nil {
		return "", fmt.Errorf("実行ファイルのパス取得に失敗: %w", err)
	}

	exeDir := filepath.Dir(exePath)
	var baseDir string

	// go run で実行されているか判定
	isGoRun := strings.Contains(exePath, "go-build") || strings.Contains(exeDir, os.TempDir())

	if isGoRun {
		// 開発時 (go run): コマンドを実行した場所（カレントディレクトリ）をベースにする
		cwd, err := osGetwd()
		if err != nil {
			return "", fmt.Errorf("カレントディレクトリの取得に失敗: %w", err)
		}
		baseDir = cwd
	} else {
		// 本番時 (EXE): 実行ファイル（バイナリ）がある場所をベースにする
		baseDir = exeDir
	}

	// ベースディレクトリと引数のフォルダ名を結合して絶対パス作成
	absolutePath := filepath.Join(baseDir, folderName)

	return absolutePath, nil
}
