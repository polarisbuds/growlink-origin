package repository

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func mimeTypeFromFilename(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "image/png"
	}
}

func readImageDataURL(uploadDir, filename string) (*string, error) {
	safeName := filepath.Base(filename)
	if safeName == "" || safeName == "." || safeName == string(filepath.Separator) {
		return nil, nil
	}
	filePath := filepath.Join(uploadDir, safeName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	mime := mimeTypeFromFilename(safeName)
	encoded := base64.StdEncoding.EncodeToString(data)
	s := fmt.Sprintf("data:%s;base64,%s", mime, encoded)
	return &s, nil
}
