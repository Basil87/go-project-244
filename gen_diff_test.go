package code

import (
	"strings"
	"testing"
)

const (
	errUnexpectedGenDiff = "unexpected error: %v"
	errGotWantGenDiff    = "got %q, want %q"
	genFile1JSON         = "file1.json"
	genFile2JSON         = "file2.json"
	genFile1YAML         = "file1.yaml"
	genFile2YAML         = "file2.yaml"
)

func TestGenDiff_StylishFormat(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempJSON(t, dir, genFile1JSON, `{"a":1,"b":2}`)
	file2 := WriteTempJSON(t, dir, genFile2JSON, `{"a":1,"b":3}`)

	got, err := GenDiff(file1, file2, "stylish")
	if err != nil {
		t.Fatalf(errUnexpectedGenDiff, err)
	}
	if !strings.Contains(got, "- b: 2") || !strings.Contains(got, "+ b: 3") {
		t.Fatalf("stylish output missing expected diff lines, got: %q", got)
	}
}

func TestGenDiff_PlainFormat(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempJSON(t, dir, genFile1JSON, `{"a":1}`)
	file2 := WriteTempJSON(t, dir, genFile2JSON, `{"a":2}`)

	got, err := GenDiff(file1, file2, "plain")
	if err != nil {
		t.Fatalf(errUnexpectedGenDiff, err)
	}
	expected := "Property 'a' was updated. From 1 to 2"
	if got != expected {
		t.Fatalf(errGotWantGenDiff, got, expected)
	}
}

func TestGenDiff_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempJSON(t, dir, genFile1JSON, `{"a":1}`)
	file2 := WriteTempJSON(t, dir, genFile2JSON, `{"a":1}`)

	got, err := GenDiff(file1, file2, "json")
	if err != nil {
		t.Fatalf(errUnexpectedGenDiff, err)
	}
	if !strings.Contains(got, `"status": "unchanged"`) {
		t.Fatalf("json output missing expected status field, got: %q", got)
	}
}

func TestGenDiff_DefaultFormat(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempJSON(t, dir, genFile1JSON, `{"x":1}`)
	file2 := WriteTempJSON(t, dir, genFile2JSON, `{"x":1}`)

	got, err := GenDiff(file1, file2, "")
	if err != nil {
		t.Fatalf(errUnexpectedGenDiff, err)
	}
	if !strings.HasPrefix(got, "{") {
		t.Fatalf("default format should return stylish output, got: %q", got)
	}
}

func TestGenDiff_YAMLFiles(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempYAML(t, dir, genFile1YAML, "a: 1\n")
	file2 := WriteTempYAML(t, dir, genFile2YAML, "a: 2\n")

	got, err := GenDiff(file1, file2, "plain")
	if err != nil {
		t.Fatalf(errUnexpectedGenDiff, err)
	}
	expected := "Property 'a' was updated. From 1 to 2"
	if got != expected {
		t.Fatalf(errGotWantGenDiff, got, expected)
	}
}

func TestGenDiff_FileNotFound(t *testing.T) {
	_, err := GenDiff("nonexistent1.json", "nonexistent2.json", "stylish")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}
