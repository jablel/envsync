package envfile

import (
	"strings"
	"testing"
)

func TestFormat_BasicEntries(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "PORT", Value: "8080"},
	}
	out := Format(entries, FormatOptions{})
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME=myapp in output, got: %s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got: %s", out)
	}
}

func TestFormat_SortKeys(t *testing.T) {
	entries := []Entry{
		{Key: "Z_KEY", Value: "z"},
		{Key: "A_KEY", Value: "a"},
		{Key: "M_KEY", Value: "m"},
	}
	out := Format(entries, FormatOptions{SortKeys: true})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected A_KEY first, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected Z_KEY last, got: %s", lines[2])
	}
}

func TestFormat_MaskSecrets(t *testing.T) {
	entries := []Entry{
		{Key: "DB_PASSWORD", Value: "supersecret"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	out := Format(entries, FormatOptions{MaskSecrets: true})
	if strings.Contains(out, "supersecret") {
		t.Errorf("expected password to be masked, got: %s", out)
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME unmasked, got: %s", out)
	}
}

func TestFormat_QuotesValuesWithSpaces(t *testing.T) {
	entries := []Entry{
		{Key: "GREETING", Value: "hello world"},
	}
	out := Format(entries, FormatOptions{})
	if !strings.Contains(out, `"hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestFormatDiff_ShowsChanges(t *testing.T) {
	result := DiffResult{
		Added:   []Entry{{Key: "NEW_KEY", Value: "newval"}},
		Removed: []Entry{{Key: "OLD_KEY", Value: "oldval"}},
		Modified: []Change{
			{Key: "CHANGED", OldValue: "v1", NewValue: "v2"},
		},
	}
	out := FormatDiff(result, false)
	if !strings.Contains(out, "+ NEW_KEY=newval") {
		t.Errorf("expected added key line, got: %s", out)
	}
	if !strings.Contains(out, "- OLD_KEY=oldval") {
		t.Errorf("expected removed key line, got: %s", out)
	}
	if !strings.Contains(out, "~ CHANGED: v1 -> v2") {
		t.Errorf("expected modified key line, got: %s", out)
	}
}

func TestFormatDiff_MasksSecrets(t *testing.T) {
	result := DiffResult{
		Modified: []Change{
			{Key: "API_SECRET", OldValue: "oldtoken", NewValue: "newtoken"},
		},
	}
	out := FormatDiff(result, true)
	if strings.Contains(out, "oldtoken") || strings.Contains(out, "newtoken") {
		t.Errorf("expected secret values to be masked, got: %s", out)
	}
}
