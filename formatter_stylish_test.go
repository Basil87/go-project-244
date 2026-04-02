package code

import (
	"testing"
)

func TestFormatStylish_Empty(t *testing.T) {
	got := FormatStylish([]diffNode{})
	expected := "{\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestFormatStylish_FlatDiff(t *testing.T) {
	nodes := []diffNode{
		{key: "a", status: statusRemoved, oldVal: float64(1)},
		{key: "a", status: statusAdded, newVal: float64(2)},
	}
	got := FormatStylish(nodes)
	expected := "{\n..- a: 1\n..+ a: 2\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestFormatStylish_Nested(t *testing.T) {
	nodes := []diffNode{
		{
			key:    "obj",
			status: statusNested,
			children: []diffNode{
				{key: "x", status: statusUnchanged, oldVal: "hello"},
			},
		},
	}
	got := FormatStylish(nodes)
	expected := "{\n..  obj: {\n......  x: hello\n....}\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestGetDiffWithFormatter_CustomFormatter(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempJSON(t, dir, "file1.json", `{"a":1}`)
	file2 := WriteTempJSON(t, dir, "file2.json", `{"a":2}`)

	customFmt := func(nodes []diffNode) string {
		return "custom"
	}
	got, err := GetDiffWithFormatter(file1, file2, customFmt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "custom" {
		t.Fatalf("got %q, want %q", got, "custom")
	}
}

func TestGetDiffWithFormatter_DefaultIsStylish(t *testing.T) {
	dir := t.TempDir()
	file1 := WriteTempJSON(t, dir, "file1.json", `{"a":1}`)
	file2 := WriteTempJSON(t, dir, "file2.json", `{"a":2}`)

	got, err := GetDiffWithFormatter(file1, file2, FormatStylish)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "{\n..- a: 1\n..+ a: 2\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}
