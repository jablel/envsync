package envfile

import (
	"testing"
)

func TestStrip_RemoveSecrets(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "API_SECRET", Value: "topsecret"},
	}
	result := Strip(entries, StripOptions{RemoveSecrets: true})
	if len(result.Kept) != 1 {
		t.Fatalf("expected 1 kept, got %d", len(result.Kept))
	}
	if result.Kept[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME kept, got %s", result.Kept[0].Key)
	}
	if len(result.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(result.Removed))
	}
}

func TestStrip_RemoveEmpty(t *testing.T) {
	entries := []Entry{
		{Key: "PRESENT", Value: "yes"},
		{Key: "EMPTY", Value: ""},
		{Key: "ALSO_EMPTY", Value: ""},
	}
	result := Strip(entries, StripOptions{RemoveEmpty: true})
	if len(result.Kept) != 1 {
		t.Fatalf("expected 1 kept, got %d", len(result.Kept))
	}
	if len(result.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(result.Removed))
	}
}

func TestStrip_RemovePrefixes(t *testing.T) {
	entries := []Entry{
		{Key: "LEGACY_HOST", Value: "old.host"},
		{Key: "LEGACY_PORT", Value: "8080"},
		{Key: "APP_HOST", Value: "new.host"},
	}
	result := Strip(entries, StripOptions{RemovePrefixes: []string{"LEGACY_"}})
	if len(result.Kept) != 1 {
		t.Fatalf("expected 1 kept, got %d", len(result.Kept))
	}
	if result.Kept[0].Key != "APP_HOST" {
		t.Errorf("expected APP_HOST kept, got %s", result.Kept[0].Key)
	}
}

func TestStrip_RemoveExplicitKeys(t *testing.T) {
	entries := []Entry{
		{Key: "KEEP_ME", Value: "yes"},
		{Key: "DROP_ME", Value: "no"},
		{Key: "ALSO_DROP", Value: "no"},
	}
	result := Strip(entries, StripOptions{RemoveKeys: []string{"DROP_ME", "ALSO_DROP"}})
	if len(result.Kept) != 1 || result.Kept[0].Key != "KEEP_ME" {
		t.Errorf("unexpected kept entries: %+v", result.Kept)
	}
	if len(result.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(result.Removed))
	}
}

func TestStrip_NoOptions_KeepsAll(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: ""},
		{Key: "SECRET_KEY", Value: "x"},
	}
	result := Strip(entries, StripOptions{})
	if len(result.Kept) != 3 {
		t.Errorf("expected all 3 kept, got %d", len(result.Kept))
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(result.Removed))
	}
}

func TestStrippedKeys(t *testing.T) {
	entries := []Entry{
		{Key: "KEEP", Value: "v"},
		{Key: "REMOVE", Value: ""},
	}
	result := Strip(entries, StripOptions{RemoveEmpty: true})
	keys := StrippedKeys(result)
	if len(keys) != 1 || keys[0] != "REMOVE" {
		t.Errorf("unexpected stripped keys: %v", keys)
	}
}
