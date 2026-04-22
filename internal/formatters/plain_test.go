package formatters

import (
	"code/internal/diff"
	"testing"
)

func TestFormatPlain_Empty(t *testing.T) {
	got := FormatPlain([]diff.DiffNode{})
	expected := ""
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatPlain_Changed(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusChanged, OldVal: "old", NewVal: "new"},
	}
	got := FormatPlain(nodes)
	expected := "Property 'a' was updated. From 'old' to 'new'"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatPlain_Nested(t *testing.T) {
	nodes := []diff.DiffNode{
		{
			Key:    "group",
			Status: diff.StatusNested,
			Children: []diff.DiffNode{
				{Key: "x", Status: diff.StatusRemoved, OldVal: float64(1)},
			},
		},
	}
	got := FormatPlain(nodes)
	expected := "Property 'group.x' was removed"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatPlain_NilValue(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusAdded, NewVal: nil},
	}
	got := FormatPlain(nodes)
	expected := "Property 'a' was added with value: null"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatPlain_ComplexValue(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusAdded, NewVal: map[string]any{"x": float64(1)}},
	}
	got := FormatPlain(nodes)
	expected := "Property 'a' was added with value: [complex value]"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatPlain_NonIntegerFloat(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "pi", Status: diff.StatusAdded, NewVal: 3.14},
	}
	got := FormatPlain(nodes)
	expected := "Property 'pi' was added with value: 3.14"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatPlain_BoolValue(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "flag", Status: diff.StatusAdded, NewVal: false},
	}
	got := FormatPlain(nodes)
	expected := "Property 'flag' was added with value: false"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatPlain_FlatDiff(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusRemoved, OldVal: float64(1)},
		{Key: "a", Status: diff.StatusAdded, NewVal: float64(2)},
	}
	got := FormatPlain(nodes)
	expected := "Property 'a' was removed\nProperty 'a' was added with value: 2"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}
