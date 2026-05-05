package envfile

import (
	"strings"
	"testing"
)

func TestFormatDiffSummary_NoChanges(t *testing.T) {
	s := SummarizeDiff([]DiffEntry{})
	got := FormatDiffSummary(s)
	if got != "no differences found" {
		t.Errorf("expected 'no differences found', got %q", got)
	}
}

func TestFormatDiffSummary_WithChanges(t *testing.T) {
	entries := []DiffEntry{
		{Key: "A", Status: StatusAdded, NewValue: "1"},
		{Key: "B", Status: StatusRemoved, OldValue: "2"},
		{Key: "C", Status: StatusModified, OldValue: "3", NewValue: "4"},
	}
	s := SummarizeDiff(entries)
	if s.Added != 1 || s.Removed != 1 || s.Modified != 1 || s.Total != 3 {
		t.Errorf("unexpected summary counts: %+v", s)
	}
	got := FormatDiffSummary(s)
	if !strings.Contains(got, "1 added") || !strings.Contains(got, "1 removed") || !strings.Contains(got, "1 modified") {
		t.Errorf("unexpected summary string: %q", got)
	}
}

func TestFormatDiffEntries_Empty(t *testing.T) {
	got := FormatDiffEntries([]DiffEntry{}, nil)
	if got != "(no changes)" {
		t.Errorf("expected '(no changes)', got %q", got)
	}
}

func TestFormatDiffEntries_ShowsSymbols(t *testing.T) {
	entries := []DiffEntry{
		{Key: "NEW_KEY", Status: StatusAdded, NewValue: "hello"},
		{Key: "OLD_KEY", Status: StatusRemoved, OldValue: "bye"},
		{Key: "MOD_KEY", Status: StatusModified, OldValue: "old", NewValue: "new"},
	}
	got := FormatDiffEntries(entries, nil)
	if !strings.Contains(got, "+ NEW_KEY=hello") {
		t.Errorf("missing added line in:\n%s", got)
	}
	if !strings.Contains(got, "- OLD_KEY=bye") {
		t.Errorf("missing removed line in:\n%s", got)
	}
	if !strings.Contains(got, "~ MOD_KEY: old -> new") {
		t.Errorf("missing modified line in:\n%s", got)
	}
}

func TestFormatDiffEntries_MasksSensitiveValues(t *testing.T) {
	masker := NewMasker(nil)
	entries := []DiffEntry{
		{Key: "SECRET_TOKEN", Status: StatusAdded, NewValue: "supersecret"},
		{Key: "DB_PASSWORD", Status: StatusModified, OldValue: "oldpass", NewValue: "newpass"},
	}
	got := FormatDiffEntries(entries, masker)
	if strings.Contains(got, "supersecret") {
		t.Errorf("sensitive value should be masked, got:\n%s", got)
	}
	if strings.Contains(got, "oldpass") || strings.Contains(got, "newpass") {
		t.Errorf("sensitive values should be masked, got:\n%s", got)
	}
}

func TestSummarizeDiff_OnlyUnchanged(t *testing.T) {
	entries := []DiffEntry{
		{Key: "A", Status: StatusUnchanged},
		{Key: "B", Status: StatusUnchanged},
	}
	s := SummarizeDiff(entries)
	if s.Total != 0 {
		t.Errorf("expected Total=0 for unchanged entries, got %d", s.Total)
	}
	got := FormatDiffSummary(s)
	if got != "no differences found" {
		t.Errorf("expected 'no differences found', got %q", got)
	}
}
