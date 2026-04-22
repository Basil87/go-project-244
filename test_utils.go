package code

import (
	"os"
	"path/filepath"
	"testing"
)

// WriteTempJSON creates a temporary JSON file in dir with the given name and content.
func WriteTempJSON(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	return path
}

// WriteTempYAML creates a temporary YAML file in dir with the given name and content.
func WriteTempYAML(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	return path
}
