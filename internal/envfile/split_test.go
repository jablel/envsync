package envfile

import (
	"testing"
)

var splitEntries = []Entry{
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "DB_PORT", Value: "5432"},
	{Key: "REDIS_URL", Value: "redis://localhost"},
	{Key: "APP_NAME", Value: "envsync"},
	{Key: "APP_ENV", Value: "production"},
	{Key: "SECRET_KEY", Value: "abc123"},
}

func TestSplit_ByExplicitKeys(t *testing.T) {
	result, err := Split(splitEntries, SplitOptions{
		Keys: map[string][]string{
			"db": {"DB_HOST", "DB_PORT"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(result.Groups["db"]); got != 2 {
		t.Errorf("expected 2 db entries, got %d", got)
	}
	if got := len(result.Unmatched); got != 4 {
		t.Errorf("expected 4 unmatched entries, got %d", got)
	}
}

func TestSplit_ByPrefix(t *testing.T) {
	result, err := Split(splitEntries, SplitOptions{
		Prefixes: map[string][]string{
			"app": {"APP_"},
			"db":  {"DB_"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(result.Groups["app"]); got != 2 {
		t.Errorf("expected 2 app entries, got %d", got)
	}
	if got := len(result.Groups["db"]); got != 2 {
		t.Errorf("expected 2 db entries, got %d", got)
	}
	if got := len(result.Unmatched); got != 2 {
		t.Errorf("expected 2 unmatched entries, got %d", got)
	}
}

func TestSplit_EntryMatchesMultipleGroups(t *testing.T) {
	result, err := Split(splitEntries, SplitOptions{
		Keys: map[string][]string{
			"all":      {"DB_HOST"},
			"database": {"DB_HOST"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(result.Groups["all"]); got != 1 {
		t.Errorf("expected 1 in 'all', got %d", got)
	}
	if got := len(result.Groups["database"]); got != 1 {
		t.Errorf("expected 1 in 'database', got %d", got)
	}
}

func TestSplit_NoOptions_ReturnsError(t *testing.T) {
	_, err := Split(splitEntries, SplitOptions{})
	if err == nil {
		t.Fatal("expected error for empty options, got nil")
	}
}

func TestSplitGroupNames_Sorted(t *testing.T) {
	result, _ := Split(splitEntries, SplitOptions{
		Prefixes: map[string][]string{
			"zebra": {"APP_"},
			"alpha": {"DB_"},
		},
	})
	names := SplitGroupNames(result)
	if len(names) != 2 {
		t.Fatalf("expected 2 group names, got %d", len(names))
	}
	if names[0] != "alpha" || names[1] != "zebra" {
		t.Errorf("expected sorted names [alpha zebra], got %v", names)
	}
}

func TestSplit_UnmatchedEmpty_WhenAllMatch(t *testing.T) {
	result, err := Split(splitEntries, SplitOptions{
		Prefixes: map[string][]string{
			"all": {"DB_", "REDIS_", "APP_", "SECRET_"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(result.Unmatched); got != 0 {
		t.Errorf("expected 0 unmatched, got %d", got)
	}
}
