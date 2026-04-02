package code

import (
	"os"
	"path/filepath"
	"testing"
)

func WriteTempJSON(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	return path
}

func WriteTempYAML(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	return path
}
