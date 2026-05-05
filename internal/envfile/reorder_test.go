package envfile

import (
	"strings"
	"testing"
)

var reorderEntries = []Entry{
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "APP_PORT", Value: "8080"},
	{Key: "SECRET_KEY", Value: "abc123"},
	{Key: "LOG_LEVEL", Value: "info"},
}

func TestReorder_ExplicitOrder(t *testing.T) {
	opts := ReorderOptions{
		Keys:          []string{"LOG_LEVEL", "APP_PORT", "DB_HOST", "SECRET_KEY"},
		UnlistedAtEnd: true,
	}
	res, err := Reorder(reorderEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Key != "LOG_LEVEL" {
		t.Errorf("expected LOG_LEVEL first, got %s", res.Entries[0].Key)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", res.Missing)
	}
}

func TestReorder_UnlistedDropped(t *testing.T) {
	opts := ReorderOptions{
		Keys:          []string{"APP_PORT", "DB_HOST"},
		UnlistedAtEnd: false,
	}
	res, err := Reorder(reorderEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
}

func TestReorder_UnlistedAtEnd(t *testing.T) {
	opts := ReorderOptions{
		Keys:          []string{"SECRET_KEY"},
		UnlistedAtEnd: true,
	}
	res, err := Reorder(reorderEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != len(reorderEntries) {
		t.Errorf("expected all entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Key != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY first, got %s", res.Entries[0].Key)
	}
}

func TestReorder_MissingKeys(t *testing.T) {
	opts := ReorderOptions{
		Keys:          []string{"NONEXISTENT", "APP_PORT"},
		UnlistedAtEnd: true,
	}
	res, err := Reorder(reorderEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "NONEXISTENT" {
		t.Errorf("expected NONEXISTENT in missing, got %v", res.Missing)
	}
}

func TestReorder_EmptyKeys_ReturnsError(t *testing.T) {
	_, err := Reorder(reorderEntries, ReorderOptions{})
	if err == nil {
		t.Error("expected error for empty key order")
	}
}

func TestReorderSummary(t *testing.T) {
	res := ReorderResult{
		Entries: reorderEntries,
		Moved:   []string{"LOG_LEVEL"},
		Missing: []string{"GONE"},
	}
	s := ReorderSummary(res)
	if !strings.Contains(s, "moved") || !strings.Contains(s, "missing") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestFormatReorderDiff_MarksMoved(t *testing.T) {
	reordered := []Entry{
		{Key: "LOG_LEVEL", Value: "info"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
	out := FormatReorderDiff(reorderEntries, reordered)
	if !strings.Contains(out, "~") {
		t.Errorf("expected moved marker in diff output: %s", out)
	}
}
