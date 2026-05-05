package envfile

import (
	"testing"
)

func freezeEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret", Comment: "# db creds"},
		{Key: "API_KEY", Value: "key123"},
	}
}

func TestFreeze_AllKeys(t *testing.T) {
	entries := freezeEntries()
	out, res, err := Freeze(entries, FreezeOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Frozen) != 3 {
		t.Errorf("expected 3 frozen, got %d", len(res.Frozen))
	}
	for _, e := range out {
		if !IsFrozen(e) {
			t.Errorf("expected %q to be frozen", e.Key)
		}
	}
}

func TestFreeze_SpecificKeys(t *testing.T) {
	entries := freezeEntries()
	out, res, err := Freeze(entries, FreezeOptions{Keys: []string{"API_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Frozen) != 1 || res.Frozen[0] != "API_KEY" {
		t.Errorf("expected only API_KEY frozen, got %v", res.Frozen)
	}
	if IsFrozen(out[0]) {
		t.Error("APP_NAME should not be frozen")
	}
	if !IsFrozen(out[2]) {
		t.Error("API_KEY should be frozen")
	}
}

func TestFreeze_SkipsAlreadyFrozen(t *testing.T) {
	entries := []Entry{
		{Key: "TOKEN", Value: "abc", Comment: "# @frozen"},
	}
	_, res, err := Freeze(entries, FreezeOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "TOKEN" {
		t.Errorf("expected TOKEN skipped, got %v", res.Skipped)
	}
	if len(res.Frozen) != 0 {
		t.Errorf("expected 0 newly frozen, got %v", res.Frozen)
	}
}

func TestFreeze_AllowUnfreeze_RefreezesEntry(t *testing.T) {
	entries := []Entry{
		{Key: "TOKEN", Value: "abc", Comment: "# @frozen"},
	}
	_, res, err := Freeze(entries, FreezeOptions{AllowUnfreeze: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Frozen) != 1 {
		t.Errorf("expected 1 frozen with AllowUnfreeze, got %v", res.Frozen)
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %v", res.Skipped)
	}
}

func TestFrozenKeys_ReturnsFrozenOnly(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1", Comment: "# @frozen"},
		{Key: "B", Value: "2"},
		{Key: "C", Value: "3", Comment: "# @frozen"},
	}
	keys := FrozenKeys(entries)
	if len(keys) != 2 {
		t.Fatalf("expected 2 frozen keys, got %d", len(keys))
	}
	if keys[0] != "A" || keys[1] != "C" {
		t.Errorf("unexpected frozen keys: %v", keys)
	}
}

func TestIsFrozen_FalseForRegularEntry(t *testing.T) {
	e := Entry{Key: "PLAIN", Value: "value", Comment: "# some comment"}
	if IsFrozen(e) {
		t.Error("expected IsFrozen to return false for regular entry")
	}
}
