package formatters

import (
	"code/diff"
	"testing"
)

func TestFormatJSON_Empty(t *testing.T) {
	got := FormatJSON([]diff.DiffNode{})
	expected := "[]"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatJSON_Added(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusAdded, NewVal: float64(2)},
	}
	got := FormatJSON(nodes)
	expected := "[\n    {\n        \"key\": \"a\",\n        \"status\": \"added\",\n        \"newValue\": 2\n    }\n]"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatJSON_Removed(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusRemoved, OldVal: float64(1)},
	}
	got := FormatJSON(nodes)
	expected := "[\n    {\n        \"key\": \"a\",\n        \"status\": \"removed\",\n        \"oldValue\": 1\n    }\n]"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatJSON_Changed(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusChanged, OldVal: float64(1), NewVal: float64(2)},
	}
	got := FormatJSON(nodes)
	expected := "[\n    {\n        \"key\": \"a\",\n        \"status\": \"changed\",\n        \"oldValue\": 1,\n        \"newValue\": 2\n    }\n]"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatJSON_Unchanged(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusUnchanged, OldVal: float64(1)},
	}
	got := FormatJSON(nodes)
	expected := "[\n    {\n        \"key\": \"a\",\n        \"status\": \"unchanged\",\n        \"oldValue\": 1\n    }\n]"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatJSON_Nested(t *testing.T) {
	nodes := []diff.DiffNode{
		{
			Key:    "group",
			Status: diff.StatusNested,
			Children: []diff.DiffNode{
				{Key: "x", Status: diff.StatusAdded, NewVal: float64(1)},
			},
		},
	}
	got := FormatJSON(nodes)
	expected := "[\n    {\n        \"key\": \"group\",\n        \"status\": \"nested\",\n        \"children\": [\n            {\n                \"key\": \"x\",\n                \"status\": \"added\",\n                \"newValue\": 1\n            }\n        ]\n    }\n]"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}
