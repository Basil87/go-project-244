package code

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempJSON(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	return path
}

func TestGetDiff_Success(t *testing.T) {
	dir := t.TempDir()

	file1 := writeTempJSON(t, dir, "file1.json", `{"a":1}`)
	file2 := writeTempJSON(t, dir, "file2.json", `{"a":2}`)

	got, err := GetDiff(file1, file2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "вот и результат" {
		t.Fatalf("got %q, want %q", got, "вот и результат")
	}
}

func TestGetDiff_FirstFileMissing(t *testing.T) {
	_, err := GetDiff("missing1.json", "missing2.json")

	if err == nil {
		t.Fatal("expected error")
	}
}
