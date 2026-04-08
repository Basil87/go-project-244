package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	errUnexpectedCmd = "unexpected error: %v"
	cmdFile1         = "f1.json"
	cmdFile2         = "f2.json"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestNewCmd_NoArgs(t *testing.T) {
	var buf bytes.Buffer
	cmd := newCmd(&buf)
	err := cmd.Run(context.Background(), []string{"gendiff"})
	if err != nil {
		t.Fatalf(errUnexpectedCmd, err)
	}
	if !strings.Contains(buf.String(), "path is required") {
		t.Fatalf("expected 'path is required', got: %q", buf.String())
	}
}

func TestNewCmd_StylishFormat(t *testing.T) {
	dir := t.TempDir()
	f1 := writeTempFile(t, dir, cmdFile1, `{"a":1,"b":2}`)
	f2 := writeTempFile(t, dir, cmdFile2, `{"a":1,"b":3}`)

	var buf bytes.Buffer
	cmd := newCmd(&buf)
	err := cmd.Run(context.Background(), []string{"gendiff", f1, f2})
	if err != nil {
		t.Fatalf(errUnexpectedCmd, err)
	}
	out := buf.String()
	if !strings.Contains(out, "- b: 2") || !strings.Contains(out, "+ b: 3") {
		t.Fatalf("stylish output missing expected diff lines, got: %q", out)
	}
}

func TestNewCmd_PlainFormat(t *testing.T) {
	dir := t.TempDir()
	f1 := writeTempFile(t, dir, cmdFile1, `{"a":1}`)
	f2 := writeTempFile(t, dir, cmdFile2, `{"a":2}`)

	var buf bytes.Buffer
	cmd := newCmd(&buf)
	err := cmd.Run(context.Background(), []string{"gendiff", "--format", "plain", f1, f2})
	if err != nil {
		t.Fatalf(errUnexpectedCmd, err)
	}
	if !strings.Contains(buf.String(), "Property 'a' was updated. From 1 to 2") {
		t.Fatalf("plain output missing expected text, got: %q", buf.String())
	}
}

func TestNewCmd_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	f1 := writeTempFile(t, dir, cmdFile1, `{"a":1}`)
	f2 := writeTempFile(t, dir, cmdFile2, `{"a":1}`)

	var buf bytes.Buffer
	cmd := newCmd(&buf)
	err := cmd.Run(context.Background(), []string{"gendiff", "--format", "json", f1, f2})
	if err != nil {
		t.Fatalf(errUnexpectedCmd, err)
	}
	if !strings.Contains(buf.String(), `"type": "unchanged"`) {
		t.Fatalf("json output missing expected field, got: %q", buf.String())
	}
}

func TestNewCmd_FormatShortFlag(t *testing.T) {
	dir := t.TempDir()
	f1 := writeTempFile(t, dir, cmdFile1, `{"x":1}`)
	f2 := writeTempFile(t, dir, cmdFile2, `{"x":2}`)

	var buf bytes.Buffer
	cmd := newCmd(&buf)
	err := cmd.Run(context.Background(), []string{"gendiff", "-f", "plain", f1, f2})
	if err != nil {
		t.Fatalf(errUnexpectedCmd, err)
	}
	if !strings.Contains(buf.String(), "Property 'x' was updated") {
		t.Fatalf("short -f flag not working, got: %q", buf.String())
	}
}

func TestNewCmd_FileNotFound(t *testing.T) {
	var buf bytes.Buffer
	cmd := newCmd(&buf)
	err := cmd.Run(context.Background(), []string{"gendiff", "nonexistent1.json", "nonexistent2.json"})
	if err == nil {
		t.Fatal("expected error for nonexistent files, got nil")
	}
}
