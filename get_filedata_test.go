package code

import (
	"strings"
	"testing"
)

func TestGetFileData_FileNotFound(t *testing.T) {
	_, err := GetFileData("nonexistent.json")

	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "file not exists") {
		t.Fatalf("expected 'file not exists' in error, got %q", err.Error())
	}
}

func TestGetFileData_IsDirectory(t *testing.T) {
	dir := t.TempDir()

	_, err := GetFileData(dir)

	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "expected file, not a directory") {
		t.Fatalf("expected 'expected file, not a directory' in error, got %q", err.Error())
	}
}

func TestGetFileData_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := WriteTempJSON(t, dir, "bad.json", `not valid json`)

	_, err := GetFileData(path)

	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid json") {
		t.Fatalf("expected 'invalid json' in error, got %q", err.Error())
	}
}

func TestGetFileData_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := WriteTempJSON(t, dir, "file.json", `{"key":"value","num":42}`)

	got, err := GetFileData(path)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "value" {
		t.Fatalf("expected key=%q, got %v", "value", got["key"])
	}
	if got["num"] != float64(42) {
		t.Fatalf("expected num=%v, got %v", float64(42), got["num"])
	}
}

func TestGetFileData_ValidYAML(t *testing.T) {
	dir := t.TempDir()
	path := WriteTempYAML(t, dir, "file.yaml", "key: value\nnum: 42\n")

	got, err := GetFileData(path)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "value" {
		t.Fatalf("expected key=%q, got %v", "value", got["key"])
	}
	if got["num"] != float64(42) {
		t.Fatalf("expected num=%v, got %v", float64(42), got["num"])
	}
}

func TestGetFileData_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := WriteTempYAML(t, dir, "bad.yaml", "key: {unclosed\n")

	_, err := GetFileData(path)

	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid yaml") {
		t.Fatalf("expected 'invalid yaml' in error, got %q", err.Error())
	}
}

func TestGetFileData_ValidYMLExtension(t *testing.T) {
	dir := t.TempDir()
	path := WriteTempYAML(t, dir, "file.yml", "a: 1\nb: hello\n")

	got, err := GetFileData(path)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["a"] != float64(1) {
		t.Fatalf("expected a=%v, got %v", float64(1), got["a"])
	}
	if got["b"] != "hello" {
		t.Fatalf("expected b=%q, got %v", "hello", got["b"])
	}
}
