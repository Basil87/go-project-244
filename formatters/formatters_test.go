package formatters

import (
	"code/diff"
	"testing"
)

func TestGetFormatter_Stylish(t *testing.T) {
	f := GetFormatter("stylish")
	got := f([]diff.DiffNode{{Key: "a", Status: diff.StatusAdded, NewVal: float64(1)}})
	expected := "{\n..+ a: 1\n}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestGetFormatter_Plain(t *testing.T) {
	f := GetFormatter("plain")
	got := f([]diff.DiffNode{{Key: "a", Status: diff.StatusRemoved, OldVal: float64(1)}})
	expected := "Property 'a' was removed"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestGetFormatter_Default(t *testing.T) {
	f := GetFormatter("unknown")
	got := f([]diff.DiffNode{})
	expected := "{\n}"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}

func TestGetFormatter_JSON(t *testing.T) {
	f := GetFormatter("json")
	got := f([]diff.DiffNode{{Key: "a", Status: diff.StatusAdded, NewVal: float64(1)}})
	expected := "[\n    {\n        \"key\": \"a\",\n        \"status\": \"added\",\n        \"newValue\": 1\n    }\n]"
	if got != expected {
		t.Fatalf(assertGotWant, got, expected)
	}
}
