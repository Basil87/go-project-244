package formatters

import (
	"code/diff"
	"testing"
)

func TestFormatPlain_Empty(t *testing.T) {
	got := FormatPlain([]diff.DiffNode{})
	expected := ""
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
