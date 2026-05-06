package envfile

import (
	"testing"
)

func TestSanitize_StripControlChars(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "hello\x01world"},
		{Key: "BAR", Value: "clean"},
	}
	opts := DefaultSanitizeOptions()
	out, results, err := Sanitize(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "helloworld" {
		t.Errorf("expected control char stripped, got %q", out[0].Value)
	}
	if results[0].Changed != true {
		t.Error("expected FOO to be marked changed")
	}
	if results[1].Changed != false {
		t.Error("expected BAR to be unchanged")
	}
}

func TestSanitize_RemoveNullBytes(t *testing.T) {
	entries := []Entry{
		{Key: "SECRET", Value: "val\x00ue"},
	}
	out, _, err := Sanitize(entries, DefaultSanitizeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "value" {
		t.Errorf("expected null byte removed, got %q", out[0].Value)
	}
}

func TestSanitize_TrimQuotes(t *testing.T) {
	entries := []Entry{
		{Key: "DB_URL", Value: `"postgres://localhost"`},
	}
	opts := DefaultSanitizeOptions()
	opts.TrimQuotes = true
	out, _, err := Sanitize(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "postgres://localhost" {
		t.Errorf("expected quotes trimmed, got %q", out[0].Value)
	}
}

func TestSanitize_NormalizeKeys(t *testing.T) {
	entries := []Entry{
		{Key: "my-key", Value: "v1"},
		{Key: "another.key", Value: "v2"},
	}
	opts := DefaultSanitizeOptions()
	opts.NormalizeKeys = true
	out, _, err := Sanitize(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "MY_KEY" {
		t.Errorf("expected MY_KEY, got %q", out[0].Key)
	}
	if out[1].Key != "ANOTHER_KEY" {
		t.Errorf("expected ANOTHER_KEY, got %q", out[1].Key)
	}
}

func TestSanitize_NoChanges(t *testing.T) {
	entries := []Entry{
		{Key: "CLEAN", Value: "already fine"},
	}
	_, results, err := Sanitize(entries, DefaultSanitizeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Changed {
		t.Error("expected no change for clean entry")
	}
}

func TestSanitizedKeys_ReturnsOnlyChanged(t *testing.T) {
	results := []SanitizeResult{
		{Key: "A", Changed: true},
		{Key: "B", Changed: false},
		{Key: "C", Changed: true},
	}
	keys := SanitizedKeys(results)
	if len(keys) != 2 || keys[0] != "A" || keys[1] != "C" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestFormatSanitizeSummary_NoChanges(t *testing.T) {
	s := SummarizeSanitize([]SanitizeResult{
		{Key: "X", Changed: false},
	})
	msg := FormatSanitizeSummary(s)
	if msg != "sanitize: 1 entries checked, no changes needed" {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestFormatSanitizeSummary_WithChanges(t *testing.T) {
	s := SummarizeSanitize([]SanitizeResult{
		{Key: "X", Changed: true},
		{Key: "Y", Changed: false},
	})
	msg := FormatSanitizeSummary(s)
	if msg != "sanitize: 1/2 entries modified" {
		t.Errorf("unexpected message: %q", msg)
	}
}
