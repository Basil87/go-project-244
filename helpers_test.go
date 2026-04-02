package code

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"- key", "key"},
		{"+ key", "key"},
		{"key", "key"},
		{"  key", "key"},
	}
	for _, c := range cases {
		got := normalize(c.input)
		if got != c.expected {
			t.Fatalf("normalize(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

func TestPrefixOrder(t *testing.T) {
	cases := []struct {
		input    string
		expected int
	}{
		{"- key", 0},
		{"+ key", 2},
		{"key", 1},
	}
	for _, c := range cases {
		got := prefixOrder(c.input)
		if got != c.expected {
			t.Fatalf("prefixOrder(%q) = %d, want %d", c.input, got, c.expected)
		}
	}
}

func TestCompareJsons_BothEmpty(t *testing.T) {
	got, err := compareJsons(map[string]any{}, map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestCompareJsons_KeyOnlyInFirst(t *testing.T) {
	got, err := compareJsons(
		map[string]any{"a": float64(1)},
		map[string]any{},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n  - a: 1\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestCompareJsons_KeyOnlyInSecond(t *testing.T) {
	got, err := compareJsons(
		map[string]any{},
		map[string]any{"b": float64(2)},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n  + b: 2\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestCompareJsons_SameKeysSameValues(t *testing.T) {
	got, err := compareJsons(
		map[string]any{"x": "hello"},
		map[string]any{"x": "hello"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n    x: hello\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestCompareJsons_SameKeysDiffValues(t *testing.T) {
	got, err := compareJsons(
		map[string]any{"x": float64(1)},
		map[string]any{"x": float64(2)},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\n  - x: 1\n  + x: 2\n}"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}
