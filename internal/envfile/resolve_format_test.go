package envfile

import (
	"strings"
	"testing"
)

func TestResolveSummary_Applied(t *testing.T) {
	r := ResolveResult{
		Entries:   []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}},
		Conflicts: nil,
		Applied:   1,
	}
	summary := ResolveSummary(r)
	if !strings.Contains(summary, "1 applied") {
		t.Errorf("expected '1 applied' in summary, got: %q", summary)
	}
	if !strings.Contains(summary, "2 total") {
		t.Errorf("expected '2 total' in summary, got: %q", summary)
	}
}

func TestResolveSummary_WithConflicts(t *testing.T) {
	r := ResolveResult{
		Entries:   []Entry{{Key: "A", Value: "1"}},
		Conflicts: []string{"A"},
		Applied:   0,
	}
	summary := ResolveSummary(r)
	if !strings.Contains(summary, "conflict") {
		t.Errorf("expected 'conflict' in summary, got: %q", summary)
	}
}

func TestFormatResolvedEntries_ShowsAdded(t *testing.T) {
	base := []Entry{{Key: "A", Value: "1"}}
	r := ResolveResult{
		Entries: []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}},
	}
	out := FormatResolvedEntries(r, base)
	if !strings.Contains(out, "+ B=2") {
		t.Errorf("expected added marker, got: %q", out)
	}
	if !strings.Contains(out, "  A=1") {
		t.Errorf("expected unchanged marker, got: %q", out)
	}
}

func TestFormatResolvedEntries_ShowsModified(t *testing.T) {
	base := []Entry{{Key: "DB_HOST", Value: "localhost"}}
	r := ResolveResult{
		Entries: []Entry{{Key: "DB_HOST", Value: "prod.example.com"}},
	}
	out := FormatResolvedEntries(r, base)
	if !strings.Contains(out, "~ DB_HOST") {
		t.Errorf("expected modified marker, got: %q", out)
	}
	if !strings.Contains(out, "was: localhost") {
		t.Errorf("expected old value in output, got: %q", out)
	}
}

func TestFormatConflicts_Empty(t *testing.T) {
	out := FormatConflicts(nil)
	if !strings.Contains(out, "No conflicts") {
		t.Errorf("expected no-conflicts message, got: %q", out)
	}
}

func TestFormatConflicts_WithKeys(t *testing.T) {
	out := FormatConflicts([]string{"SECRET_KEY", "DB_PASS"})
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output, got: %q", out)
	}
	if !strings.Contains(out, "2 keys") {
		t.Errorf("expected count in output, got: %q", out)
	}
}
