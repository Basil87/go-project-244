package code

import (
	"testing"
)

func TestGetDiff_Success(t *testing.T) {
	dir := t.TempDir()

	file1 := WriteTempJSON(t, dir, "file1.json", `{"a":1}`)
	file2 := WriteTempJSON(t, dir, "file2.json", `{"a":2}`)

	got, err := GetDiff(file1, file2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n  - a: 1\n  + a: 2\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestGetDiff_FirstFileMissing(t *testing.T) {
	_, err := GetDiff("missing1.json", "missing2.json")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetDiff_SecondFileMissing(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempJSON(t, dir, "file1.json", `{"a":1}`)

	_, err := GetDiff(file1, "missing2.json")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetDiff_IdenticalFiles(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempJSON(t, dir, "file1.json", `{"a":1}`)
	file2 := WriteTempJSON(t, dir, "file2.json", `{"a":1}`)

	got, err := GetDiff(file1, file2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n    a: 1\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestGetDiff_KeysOnlyInFirstFile(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempJSON(t, dir, "file1.json", `{"a":1}`)
	file2 := WriteTempJSON(t, dir, "file2.json", `{}`)

	got, err := GetDiff(file1, file2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n  - a: 1\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestGetDiff_KeysOnlyInSecondFile(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempJSON(t, dir, "file1.json", `{}`)
	file2 := WriteTempJSON(t, dir, "file2.json", `{"b":2}`)

	got, err := GetDiff(file1, file2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n  + b: 2\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestGetDiff_YAMLFilesChangedValue(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempYAML(t, dir, "file1.yaml", "a: 1\n")
	file2 := WriteTempYAML(t, dir, "file2.yaml", "a: 2\n")

	got, err := GetDiff(file1, file2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n  - a: 1\n  + a: 2\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestGetDiff_YAMLFilesIdentical(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempYAML(t, dir, "file1.yaml", "host: hexlet.io\n")
	file2 := WriteTempYAML(t, dir, "file2.yaml", "host: hexlet.io\n")

	got, err := GetDiff(file1, file2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n    host: hexlet.io\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestGetDiff_YAMLFilesKeyRemoved(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempYAML(t, dir, "file1.yaml", "a: 1\nb: 2\n")
	file2 := WriteTempYAML(t, dir, "file2.yaml", "a: 1\n")

	got, err := GetDiff(file1, file2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n    a: 1\n  - b: 2\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestGetDiff_MixedJSONAndYAML(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempJSON(t, dir, "file1.json", `{"a":1}`)
	file2 := WriteTempYAML(t, dir, "file2.yaml", "a: 2\n")

	got, err := GetDiff(file1, file2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n  - a: 1\n  + a: 2\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}
