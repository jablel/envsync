package envfile

import (
	"testing"
)

func TestRedact_NonSensitiveUnchanged(t *testing.T) {
	masker := NewMasker(nil)
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "PORT", Value: "8080"},
	}
	result := Redact(entries, masker, DefaultRedactOptions())
	for i, e := range result {
		if e.Value != entries[i].Value {
			t.Errorf("key %s: expected %q, got %q", e.Key, entries[i].Value, e.Value)
		}
	}
}

func TestRedact_SensitiveReplaced(t *testing.T) {
	masker := NewMasker(nil)
	entries := []Entry{
		{Key: "DB_PASSWORD", Value: "supersecret"},
		{Key: "API_SECRET", Value: "abc123"},
	}
	result := Redact(entries, masker, DefaultRedactOptions())
	for _, e := range result {
		if e.Value != "[REDACTED]" {
			t.Errorf("key %s: expected [REDACTED], got %q", e.Key, e.Value)
		}
	}
}

func TestRedact_CustomReplacement(t *testing.T) {
	masker := NewMasker(nil)
	opts := RedactOptions{Replacement: "***"}
	entries := []Entry{
		{Key: "SECRET_KEY", Value: "topsecret"},
	}
	result := Redact(entries, masker, opts)
	if result[0].Value != "***" {
		t.Errorf("expected ***, got %q", result[0].Value)
	}
}

func TestRedact_PartialReveal(t *testing.T) {
	masker := NewMasker(nil)
	opts := RedactOptions{Replacement: "[REDACTED]", PartialReveal: 3}
	entries := []Entry{
		{Key: "API_TOKEN", Value: "abcdef"},
	}
	result := Redact(entries, masker, opts)
	// First 3 chars visible, rest replaced with '*'
	expected := "abc***"
	if result[0].Value != expected {
		t.Errorf("expected %q, got %q", expected, result[0].Value)
	}
}

func TestRedact_PartialRevealShortValue(t *testing.T) {
	masker := NewMasker(nil)
	opts := RedactOptions{Replacement: "[REDACTED]", PartialReveal: 10}
	entries := []Entry{
		{Key: "DB_PASS", Value: "hi"},
	}
	result := Redact(entries, masker, opts)
	// Value shorter than PartialReveal — fully redacted
	if result[0].Value != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", result[0].Value)
	}
}

func TestRedactMap_SensitiveKeys(t *testing.T) {
	masker := NewMasker(nil)
	m := map[string]string{
		"APP_ENV":    "production",
		"DB_PASSWORD": "s3cr3t",
	}
	result := RedactMap(m, masker, DefaultRedactOptions())
	if result["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged")
	}
	if result["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", result["DB_PASSWORD"])
	}
}

func TestRedactMap_EmptyReplacement_UsesDefault(t *testing.T) {
	masker := NewMasker(nil)
	opts := RedactOptions{Replacement: ""}
	m := map[string]string{"API_KEY": "mykey"}
	result := RedactMap(m, masker, opts)
	if result["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", result["API_KEY"])
	}
}
