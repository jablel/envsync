package envfile

import (
	"testing"
)

func TestNormalize_TrimWhitespace(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "  localhost  "},
		{Key: "PORT", Value: "8080"},
	}
	opts := DefaultNormalizeOptions()
	opts.UppercaseKeys = false
	opts.RemoveEmptyValues = false

	res := Normalize(entries, opts)

	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Value != "localhost" {
		t.Errorf("expected 'localhost', got %q", res.Entries[0].Value)
	}
	if len(res.Changed) != 1 || res.Changed[0] != "HOST" {
		t.Errorf("expected Changed=[HOST], got %v", res.Changed)
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	entries := []Entry{
		{Key: "db_host", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
	}
	opts := DefaultNormalizeOptions()
	opts.UppercaseKeys = true

	res := Normalize(entries, opts)

	if res.Entries[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", res.Entries[0].Key)
	}
	if res.Entries[1].Key != "DB_PORT" {
		t.Errorf("expected DB_PORT, got %q", res.Entries[1].Key)
	}
	if len(res.Changed) != 1 {
		t.Errorf("expected 1 changed key, got %v", res.Changed)
	}
}

func TestNormalize_RemoveEmptyValues(t *testing.T) {
	entries := []Entry{
		{Key: "PRESENT", Value: "yes"},
		{Key: "EMPTY", Value: ""},
	}
	opts := DefaultNormalizeOptions()
	opts.RemoveEmptyValues = true

	res := Normalize(entries, opts)

	if len(res.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(res.Entries))
	}
	if res.Entries[0].Key != "PRESENT" {
		t.Errorf("expected PRESENT, got %q", res.Entries[0].Key)
	}
	if len(res.Dropped) != 1 || res.Dropped[0] != "EMPTY" {
		t.Errorf("expected Dropped=[EMPTY], got %v", res.Dropped)
	}
}

func TestNormalize_NormalizeLineEndings(t *testing.T) {
	entries := []Entry{
		{Key: "MULTILINE", Value: "line1\r\nline2\rline3"},
	}
	opts := DefaultNormalizeOptions()

	res := Normalize(entries, opts)

	expected := "line1\nline2\nline3"
	if res.Entries[0].Value != expected {
		t.Errorf("expected %q, got %q", expected, res.Entries[0].Value)
	}
	if len(res.Changed) != 1 {
		t.Errorf("expected 1 changed key, got %v", res.Changed)
	}
}

func TestNormalize_NoChanges(t *testing.T) {
	entries := []Entry{
		{Key: "KEY", Value: "value"},
	}
	opts := DefaultNormalizeOptions()

	res := Normalize(entries, opts)

	if len(res.Changed) != 0 {
		t.Errorf("expected no changes, got %v", res.Changed)
	}
	if len(res.Dropped) != 0 {
		t.Errorf("expected no drops, got %v", res.Dropped)
	}
}
