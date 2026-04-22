package formatters

import (
	"code/internal/diff"
	"testing"
)

func TestFormatJSON_Empty(t *testing.T) {
	got := FormatJSON([]diff.DiffNode{})
	expected := "{}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatJSON_Added(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusAdded, NewVal: float64(2)},
	}
	got := FormatJSON(nodes)
	expected := "{\n    \"a\": {\n        \"type\": \"added\",\n        \"value\": 2\n    }\n}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatJSON_Removed(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusRemoved, OldVal: float64(1)},
	}
	got := FormatJSON(nodes)
	expected := "{\n    \"a\": {\n        \"type\": \"removed\",\n        \"value\": 1\n    }\n}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatJSON_Changed(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusChanged, OldVal: float64(1), NewVal: float64(2)},
	}
	got := FormatJSON(nodes)
	expected := "{\n    \"a\": {\n        \"from\": 1,\n        \"to\": 2,\n        \"type\": \"changed\"\n    }\n}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestFormatJSON_Unchanged(t *testing.T) {
	nodes := []diff.DiffNode{
		{Key: "a", Status: diff.StatusUnchanged, OldVal: float64(1)},
	}
	got := FormatJSON(nodes)
	expected := "{\n    \"a\": {\n        \"type\": \"unchanged\",\n        \"value\": 1\n    }\n}"
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
	expected := "{\n    \"group\": {\n        \"children\": {\n            \"x\": {\n                \"type\": \"added\",\n                \"value\": 1\n            }\n        },\n        \"type\": \"nested\"\n    }\n}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}
