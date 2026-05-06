package envfile

import (
	"testing"
)

func TestCompact_RemovesEmptyValues(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: ""},
		{Key: "DB", Value: "mydb"},
	}
	opts := DefaultCompactOptions()
	opts.RemoveEmpty = true
	res := Compact(entries, opts)
	if res.Kept != 2 {
		t.Errorf("expected 2 kept, got %d", res.Kept)
	}
	if len(res.Removed) != 1 || res.Removed[0].Key != "PORT" {
		t.Errorf("expected PORT removed, got %+v", res.Removed)
	}
}

func TestCompact_RemovesCommentedPlaceholders(t *testing.T) {
	entries := []Entry{
		{Key: "API_KEY", Value: "#changeme"},
		{Key: "HOST", Value: "prod.example.com"},
	}
	opts := DefaultCompactOptions()
	opts.RemoveCommented = true
	res := Compact(entries, opts)
	if res.Kept != 1 {
		t.Errorf("expected 1 kept, got %d", res.Kept)
	}
	if res.Output[0].Key != "HOST" {
		t.Errorf("expected HOST in output, got %s", res.Output[0].Key)
	}
}

func TestCompact_DeduplicatesKeepsLast(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "first"},
		{Key: "BAR", Value: "bar"},
		{Key: "FOO", Value: "last"},
	}
	opts := DefaultCompactOptions() // RemoveDuplicates: true
	res := Compact(entries, opts)
	if res.Kept != 2 {
		t.Errorf("expected 2 kept, got %d", res.Kept)
	}
	var fooVal string
	for _, e := range res.Output {
		if e.Key == "FOO" {
			fooVal = e.Value
		}
	}
	if fooVal != "last" {
		t.Errorf("expected FOO=last, got %s", fooVal)
	}
}

func TestCompact_RemovesDefaults(t *testing.T) {
	entries := []Entry{
		{Key: "LOG_LEVEL", Value: "info"},
		{Key: "TIMEOUT", Value: "30"},
		{Key: "HOST", Value: "localhost"},
	}
	opts := DefaultCompactOptions()
	opts.RemoveDefaults = map[string]string{
		"LOG_LEVEL": "info",
		"HOST":      "localhost",
	}
	res := Compact(entries, opts)
	if res.Kept != 1 {
		t.Errorf("expected 1 kept, got %d", res.Kept)
	}
	if res.Output[0].Key != "TIMEOUT" {
		t.Errorf("expected TIMEOUT in output, got %s", res.Output[0].Key)
	}
}

func TestCompact_NoDuplicates_Unchanged(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	opts := DefaultCompactOptions()
	res := Compact(entries, opts)
	if res.Kept != 2 {
		t.Errorf("expected 2 kept, got %d", res.Kept)
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected nothing removed, got %+v", res.Removed)
	}
}

func TestCompact_EmptyInput(t *testing.T) {
	res := Compact([]Entry{}, DefaultCompactOptions())
	if res.Kept != 0 {
		t.Errorf("expected 0 kept, got %d", res.Kept)
	}
}
