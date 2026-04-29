package envfile

import (
	"testing"
)

var pinEntries = []Entry{
	{Key: "APP_ENV", Value: "development"},
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "API_KEY", Value: "old-key"},
	{Key: "LOG_LEVEL", Value: "debug"},
}

func TestPin_UpdatesMatchingKeys(t *testing.T) {
	pins := map[string]string{"APP_ENV": "production", "LOG_LEVEL": "info"}
	out, res, err := Pin(pinEntries, pins, PinOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 2 {
		t.Errorf("expected 2 pinned, got %d", len(res.Pinned))
	}
	for _, e := range out {
		if e.Key == "APP_ENV" && e.Value != "production" {
			t.Errorf("APP_ENV not pinned to production, got %q", e.Value)
		}
		if e.Key == "LOG_LEVEL" && e.Value != "info" {
			t.Errorf("LOG_LEVEL not pinned to info, got %q", e.Value)
		}
	}
}

func TestPin_SkipsAlreadyMatchingValue(t *testing.T) {
	pins := map[string]string{"DB_HOST": "localhost"}
	_, res, err := Pin(pinEntries, pins, PinOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST skipped, got %v", res.Skipped)
	}
}

func TestPin_OverwriteForcesPinEvenIfSameValue(t *testing.T) {
	pins := map[string]string{"DB_HOST": "localhost"}
	_, res, err := Pin(pinEntries, pins, PinOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 1 || res.Pinned[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST pinned with overwrite, got %v", res.Pinned)
	}
}

func TestPin_MissingKeyRecorded(t *testing.T) {
	pins := map[string]string{"MISSING_KEY": "value"}
	_, res, err := Pin(pinEntries, pins, PinOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY in missing, got %v", res.Missing)
	}
}

func TestPin_FailOnMissingReturnsError(t *testing.T) {
	pins := map[string]string{"GHOST": "val"}
	_, _, err := Pin(pinEntries, pins, PinOptions{FailOnMissing: true})
	if err == nil {
		t.Error("expected error for missing key with FailOnMissing, got nil")
	}
}

func TestPin_OriginalEntriesUnmodified(t *testing.T) {
	orig := make([]Entry, len(pinEntries))
	copy(orig, pinEntries)
	pins := map[string]string{"APP_ENV": "staging"}
	Pin(pinEntries, pins, PinOptions{})
	for i, e := range pinEntries {
		if e != orig[i] {
			t.Errorf("original entries modified at index %d", i)
		}
	}
}

func TestPinnedKeys_ReturnsMatches(t *testing.T) {
	pins := map[string]string{"APP_ENV": "development", "DB_HOST": "remote", "LOG_LEVEL": "debug"}
	matched := PinnedKeys(pinEntries, pins)
	if len(matched) != 2 {
		t.Errorf("expected 2 matched pinned keys, got %d: %v", len(matched), matched)
	}
}
