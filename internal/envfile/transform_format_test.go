package envfile

import (
	"strings"
	"testing"
)

func TestSummarizeTransform_NoChanges(t *testing.T) {
	entries := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	s := SummarizeTransform(entries, entries)
	if s.Total != 2 || s.Changed != 0 {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestSummarizeTransform_WithChanges(t *testing.T) {
	before := []Entry{{Key: "A", Value: "old"}, {Key: "B", Value: "same"}}
	after := []Entry{{Key: "A", Value: "new"}, {Key: "B", Value: "same"}}
	s := SummarizeTransform(before, after)
	if s.Changed != 1 {
		t.Errorf("expected 1 change, got %d", s.Changed)
	}
}

func TestFormatTransformSummary_Basic(t *testing.T) {
	s := TransformSummary{Total: 5, Changed: 3}
	out := FormatTransformSummary(s)
	if !strings.Contains(out, "5 entries") {
		t.Errorf("missing total: %s", out)
	}
	if !strings.Contains(out, "3 changed") {
		t.Errorf("missing changed count: %s", out)
	}
}

func TestFormatTransformSummary_WithSkipped(t *testing.T) {
	s := TransformSummary{Total: 4, Changed: 2, Skipped: 1}
	out := FormatTransformSummary(s)
	if !strings.Contains(out, "skipped") {
		t.Errorf("expected skipped mention: %s", out)
	}
}

func TestFormatTransformedEntries_ShowsChanges(t *testing.T) {
	before := []Entry{{Key: "HOST", Value: "old"}, {Key: "PORT", Value: "80"}}
	after := []Entry{{Key: "HOST", Value: "new"}, {Key: "PORT", Value: "80"}}
	out := FormatTransformedEntries(before, after, nil)
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected HOST in output: %s", out)
	}
	if strings.Contains(out, "PORT") {
		t.Errorf("PORT should not appear (unchanged): %s", out)
	}
}

func TestFormatTransformedEntries_NoChanges(t *testing.T) {
	entries := []Entry{{Key: "A", Value: "1"}}
	out := FormatTransformedEntries(entries, entries, nil)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected no-changes message: %s", out)
	}
}

func TestFormatTransformedEntries_MasksSensitive(t *testing.T) {
	m := NewMasker(nil)
	before := []Entry{{Key: "API_SECRET", Value: "old_secret"}}
	after := []Entry{{Key: "API_SECRET", Value: "new_secret"}}
	out := FormatTransformedEntries(before, after, m)
	if strings.Contains(out, "old_secret") || strings.Contains(out, "new_secret") {
		t.Errorf("sensitive values should be masked: %s", out)
	}
	if !strings.Contains(out, "API_SECRET") {
		t.Errorf("key should still appear: %s", out)
	}
}
