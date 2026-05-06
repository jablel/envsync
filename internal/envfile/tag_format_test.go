package envfile

import (
	"strings"
	"testing"
)

func TestSummarizeTag_NoResults(t *testing.T) {
	s := SummarizeTag(nil)
	if s.TotalTagged != 0 || s.TotalSkipped != 0 {
		t.Errorf("expected zero counts, got tagged=%d skipped=%d", s.TotalTagged, s.TotalSkipped)
	}
}

func TestSummarizeTag_WithResults(t *testing.T) {
	results := []TagResult{
		{Key: "DB_HOST", Tagged: []string{"env:prod", "team:infra"}},
		{Key: "PORT", Skipped: []string{"env:prod"}},
	}
	s := SummarizeTag(results)
	if s.TotalTagged != 2 {
		t.Errorf("expected 2 tagged, got %d", s.TotalTagged)
	}
	if s.TotalSkipped != 1 {
		t.Errorf("expected 1 skipped, got %d", s.TotalSkipped)
	}
}

func TestFormatTagSummary_NoMatches(t *testing.T) {
	s := SummarizeTag(nil)
	out := FormatTagSummary(s)
	if !strings.Contains(out, "no entries matched") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatTagSummary_WithChanges(t *testing.T) {
	results := []TagResult{
		{Key: "DB_HOST", Tagged: []string{"env:prod"}},
	}
	s := SummarizeTag(results)
	out := FormatTagSummary(s)
	if !strings.Contains(out, "1 tag(s) applied") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatTaggedEntries_ShowsSymbols(t *testing.T) {
	results := []TagResult{
		{Key: "DB_HOST", Tagged: []string{"env:prod"}},
		{Key: "PORT", Skipped: []string{"env:prod"}},
	}
	out := FormatTaggedEntries(results)
	if !strings.Contains(out, "+ DB_HOST") {
		t.Errorf("expected '+' symbol for tagged entry, got: %q", out)
	}
	if !strings.Contains(out, "~ PORT") {
		t.Errorf("expected '~' symbol for skipped entry, got: %q", out)
	}
}

func TestFormatTaggedEntries_Empty(t *testing.T) {
	out := FormatTaggedEntries(nil)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected no-changes message, got: %q", out)
	}
}
