package formatters

import (
	"code/internal/diff"
	"testing"
)

const assertGotWant = "got %q, want %q"

func TestFormatStylish_Empty(t *testing.T) {
	got := FormatStylish([]diff.DiffNode{})
	expected := "{\n}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatStylish_FlatDiff(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusRemoved, OldVal: float64(1)},
		{Key: "a", Status: diff.StatusAdded, NewVal: float64(2)},
	}
	got := FormatStylish(nodes)
	expected := "{\n  - a: 1\n  + a: 2\n}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatStylish_MapValue(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "obj", Status: diff.StatusRemoved, OldVal: map[string]any{"x": float64(1)}},
	}
	got := FormatStylish(nodes)
	expected := "{\n  - obj: {\n        x: 1\n    }\n}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatStylish_NilValue(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusAdded, NewVal: nil},
	}
	got := FormatStylish(nodes)
	expected := "{\n  + a: null\n}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatStylish_NonIntegerFloat(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "pi", Status: diff.StatusAdded, NewVal: 3.14},
	}
	got := FormatStylish(nodes)
	expected := "{\n  + pi: 3.14\n}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatStylish_Nested(t *testing.T) {
	nodes := []diff.DiffNode{
		{
			Key:    "obj",
			Status: diff.StatusNested,
			Children: []diff.DiffNode{
				{Key: "x", Status: diff.StatusUnchanged, OldVal: "hello"},
			},
		},
	}
	got := FormatStylish(nodes)
	expected := "{\n    obj: {\n        x: hello\n    }\n}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}
