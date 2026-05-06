package envfile

import (
	"strings"
	"testing"
)

func TestFormatDiffSnapshotSummary_NoChanges(t *testing.T) {
	s := DiffSnapshotSummary{Unchanged: 5}
	out := FormatDiffSnapshotSummary(s)
	if !strings.Contains(out, "No changes") {
		t.Errorf("expected no-changes message, got: %s", out)
	}
	if !strings.Contains(out, "5 unchanged") {
		t.Errorf("expected unchanged count, got: %s", out)
	}
}

func TestFormatDiffSnapshotSummary_WithChanges(t *testing.T) {
	s := DiffSnapshotSummary{Added: 2, Removed: 1, Modified: 3, Unchanged: 4}
	out := FormatDiffSnapshotSummary(s)
	if !strings.Contains(out, "+2 added") {
		t.Errorf("expected added count, got: %s", out)
	}
	if !strings.Contains(out, "-1 removed") {
		t.Errorf("expected removed count, got: %s", out)
	}
	if !strings.Contains(out, "~3 modified") {
		t.Errorf("expected modified count, got: %s", out)
	}
}

func TestSummarizeDiffSnapshot_Counts(t *testing.T) {
	d := DiffResult{
		Added:     []Entry{{Key: "A"}},
		Removed:   []Entry{{Key: "B"}, {Key: "C"}},
		Modified:  []Change{{Old: Entry{Key: "D"}, New: Entry{Key: "D"}}},
		Unchanged: []Entry{{Key: "E"}, {Key: "F"}, {Key: "G"}},
	}
	s := SummarizeDiffSnapshot(d)
	if s.Added != 1 || s.Removed != 2 || s.Modified != 1 || s.Unchanged != 3 {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestFormatDiffSnapshotEntries_ShowsSymbols(t *testing.T) {
	d := DiffResult{
		Added:   []Entry{{Key: "NEW", Value: "val"}},
		Removed: []Entry{{Key: "OLD", Value: "gone"}},
		Modified: []Change{
			{Old: Entry{Key: "MOD", Value: "before"}, New: Entry{Key: "MOD", Value: "after"}},
		},
	}
	out := FormatDiffSnapshotEntries(d, nil)
	if !strings.Contains(out, "+ NEW=val") {
		t.Errorf("expected added line, got: %s", out)
	}
	if !strings.Contains(out, "- OLD=gone") {
		t.Errorf("expected removed line, got: %s", out)
	}
	if !strings.Contains(out, "~ MOD: before -> after") {
		t.Errorf("expected modified line, got: %s", out)
	}
}

func TestFormatDiffSnapshotEntries_MasksSensitiveValues(t *testing.T) {
	m := NewMasker(nil)
	d := DiffResult{
		Added: []Entry{{Key: "SECRET_TOKEN", Value: "supersecret"}},
	}
	out := FormatDiffSnapshotEntries(d, m)
	if strings.Contains(out, "supersecret") {
		t.Errorf("sensitive value should be masked, got: %s", out)
	}
	if !strings.Contains(out, "SECRET_TOKEN") {
		t.Errorf("key should still appear, got: %s", out)
	}
}
