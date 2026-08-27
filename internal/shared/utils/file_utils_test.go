package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetAbsoluteTargetDir(t *testing.T) {
	origExecutable := osExecutable
	origGetwd := osGetwd
	t.Cleanup(func() {
		osExecutable = origExecutable
		osGetwd = origGetwd
	})

	exeBaseDir := filepath.Join(string(filepath.Separator), "opt", "growlink")
	exePath := filepath.Join(exeBaseDir, "growlink")
	cwd := filepath.Join(string(filepath.Separator), "workspace", "growlink")
	goRunPath := filepath.Join(os.TempDir(), "go-build12345", "exe")
	tempExePath := filepath.Join(os.TempDir(), "growlink-test-bin")

	type args struct {
		folderName string
	}
	tests := []struct {
		name        string
		args        args
		executable  func() (string, error)
		getwd       func() (string, error)
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name: "通常系: EXE相当（実行ファイルのディレクトリがベース）",
			args: args{folderName: "example"},
			executable: func() (string, error) {
				return exePath, nil
			},
			getwd: func() (string, error) {
				t.Fatal("EXE相当では os.Getwd は呼ばれない")
				return "", nil
			},
			want: filepath.Join(exeBaseDir, "example"),
		},
		{
			name: "通常系: go run 相当（パスに go-build を含み、カレントディレクトリがベース）",
			args: args{folderName: "example"},
			executable: func() (string, error) {
				return goRunPath, nil
			},
			getwd: func() (string, error) {
				return cwd, nil
			},
			want: filepath.Join(cwd, "example"),
		},
		{
			name: "通常系: go run 相当（実行ディレクトリが TempDir 配下で、カレントディレクトリがベース）",
			args: args{folderName: "example"},
			executable: func() (string, error) {
				return tempExePath, nil
			},
			getwd: func() (string, error) {
				return cwd, nil
			},
			want: filepath.Join(cwd, "example"),
		},
		{
			name: "異常系: 実行ファイルのパス取得に失敗",
			args: args{folderName: "example"},
			executable: func() (string, error) {
				return "", errors.New("executable not found")
			},
			getwd: func() (string, error) {
				t.Fatal("os.Executable 失敗時は os.Getwd は呼ばれない")
				return "", nil
			},
			wantErr:     true,
			errContains: "実行ファイルのパス取得に失敗",
		},
		{
			name: "異常系: go run 相当でカレントディレクトリの取得に失敗",
			args: args{folderName: "example"},
			executable: func() (string, error) {
				return goRunPath, nil
			},
			getwd: func() (string, error) {
				return "", errors.New("getcwd failed")
			},
			wantErr:     true,
			errContains: "カレントディレクトリの取得に失敗",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osExecutable = tt.executable
			osGetwd = tt.getwd

			got, err := GetAbsoluteTargetDir(tt.args.folderName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAbsoluteTargetDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetAbsoluteTargetDir() error = %v, want contains %q", err, tt.errContains)
				}
				if got != "" {
					t.Errorf("GetAbsoluteTargetDir() = %v, want empty on error", got)
				}
				return
			}
			if got != tt.want {
				t.Errorf("GetAbsoluteTargetDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
