package envfile

import (
	"strings"
	"testing"
)

func TestRotateSummary_WithRotatedKeys(t *testing.T) {
	result := RotateResult{
		Rotated: []string{"DB_PASSWORD", "API_SECRET"},
		Skipped: []string{"APP_NAME", "PORT"},
		Errors:  map[string]error{},
	}
	summary := RotateSummary(result)
	if !strings.Contains(summary, "Rotated: 2") {
		t.Errorf("expected rotated count, got: %s", summary)
	}
	if !strings.Contains(summary, "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD in summary")
	}
	if !strings.Contains(summary, "Skipped: 2") {
		t.Errorf("expected skipped count, got: %s", summary)
	}
}

func TestRotateSummary_WithErrors(t *testing.T) {
	result := RotateResult{
		Rotated: []string{},
		Skipped: []string{"APP_NAME"},
		Errors:  map[string]error{"DB_PASSWORD": errorf("vault down")},
	}
	summary := RotateSummary(result)
	if !strings.Contains(summary, "Failed:  1") {
		t.Errorf("expected failed count, got: %s", summary)
	}
	if !strings.Contains(summary, "DB_PASSWORD") {
		t.Error("expected DB_PASSWORD in error section")
	}
}

func TestFormatRotateInstructions_MasksValues(t *testing.T) {
	original := []Entry{
		{Key: "DB_PASSWORD", Value: "old-pass"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	updated := []Entry{
		{Key: "DB_PASSWORD", Value: "new-pass"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	masker := NewMasker(nil)
	out := FormatRotateInstructions(original, updated, masker)

	if !strings.Contains(out, "rotate DB_PASSWORD") {
		t.Errorf("expected rotate instruction for DB_PASSWORD, got: %s", out)
	}
	if strings.Contains(out, "old-pass") || strings.Contains(out, "new-pass") {
		t.Error("sensitive values should be masked in instructions")
	}
	if strings.Contains(out, "APP_NAME") {
		t.Error("unchanged key APP_NAME should not appear")
	}
}

func TestFormatRotateInstructions_NoChanges(t *testing.T) {
	entries := []Entry{{Key: "PORT", Value: "8080"}}
	masker := NewMasker(nil)
	out := FormatRotateInstructions(entries, entries, masker)
	if out != "" {
		t.Errorf("expected empty output for no changes, got: %s", out)
	}
}

// errorf is a local helper to create errors in tests without importing errors package.
func errorf(msg string) error {
	return &rotateTestError{msg}
}

type rotateTestError struct{ msg string }

func (e *rotateTestError) Error() string { return e.msg }
