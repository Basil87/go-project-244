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

	expected := "{\n..- a: 1\n..+ a: 2\n}"
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

	expected := "{\n..  a: 1\n}"
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

	expected := "{\n..- a: 1\n}"
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

	expected := "{\n..+ b: 2\n}"
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

	expected := "{\n..- a: 1\n..+ a: 2\n}"
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

	expected := "{\n..  host: hexlet.io\n}"
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

	expected := "{\n..  a: 1\n..- b: 2\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestGetDiff_NestedStructures(t *testing.T) {
	got, err := GetDiff("mock/file1.json", "mock/file2.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n" +
		"..  common: {\n" +
		"......+ follow: false\n" +
		"......  setting1: Value 1\n" +
		"......- setting2: 200\n" +
		"......- setting3: true\n" +
		"......+ setting3: null\n" +
		"......+ setting4: blah blah\n" +
		"......+ setting5: {\n" +
		"..........  key5: value5\n" +
		"........}\n" +
		"......  setting6: {\n" +
		"..........  doge: {\n" +
		"..............- wow: \n" +
		"..............+ wow: so much\n" +
		"............}\n" +
		"..........  key: value\n" +
		"..........+ ops: vops\n" +
		"........}\n" +
		"....}\n" +
		"..  group1: {\n" +
		"......- baz: bas\n" +
		"......+ baz: bars\n" +
		"......  foo: bar\n" +
		"......- nest: {\n" +
		"..........  key: value\n" +
		"........}\n" +
		"......+ nest: str\n" +
		"....}\n" +
		"..- group2: {\n" +
		"......  abc: 12345\n" +
		"......  deep: {\n" +
		"..........  id: 45\n" +
		"........}\n" +
		"....}\n" +
		"..+ group3: {\n" +
		"......  deep: {\n" +
		"..........  id: {\n" +
		"..............  number: 45\n" +
		"............}\n" +
		"........}\n" +
		"......  fee: 100500\n" +
		"....}\n" +
		"}"
	if got != expected {
		t.Fatalf("got:\n%s\n\nwant:\n%s", got, expected)
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

	expected := "{\n..- a: 1\n..+ a: 2\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}
