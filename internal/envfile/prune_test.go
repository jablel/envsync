package envfile

import (
	"testing"
)

func pruneEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "CHANGEME"},
		{Key: "DEBUG_FLAG", Value: "true"},
		{Key: "API_KEY", Value: "your_api_key_here"},
		{Key: "PORT", Value: "8080"},
		{Key: "LEGACY_HOST", Value: "old.example.com"},
		{Key: "COMMENTED_VAL", Value: "# not set"},
	}
}

func TestPrune_RemoveExplicitKeys(t *testing.T) {
	result := Prune(pruneEntries(), PruneOptions{
		RemoveKeys: []string{"DEBUG_FLAG", "PORT"},
	})
	if len(result.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(result.Removed))
	}
	for _, e := range result.Removed {
		if e.Key != "DEBUG_FLAG" && e.Key != "PORT" {
			t.Errorf("unexpected removed key: %s", e.Key)
		}
	}
}

func TestPrune_RemovePlaceholders(t *testing.T) {
	result := Prune(pruneEntries(), PruneOptions{
		RemovePlaceholders: true,
	})
	removedKeys := PrunedKeys(result)
	if !contains(removedKeys, "DB_PASSWORD") {
		t.Error("expected DB_PASSWORD (CHANGEME) to be pruned")
	}
	if !contains(removedKeys, "API_KEY") {
		t.Error("expected API_KEY (your_api_key_here) to be pruned")
	}
	if contains(removedKeys, "APP_NAME") {
		t.Error("APP_NAME should not be pruned")
	}
}

func TestPrune_RemoveByPrefix(t *testing.T) {
	result := Prune(pruneEntries(), PruneOptions{
		RemovePrefixes: []string{"LEGACY_", "DEBUG_"},
	})
	removedKeys := PrunedKeys(result)
	if !contains(removedKeys, "LEGACY_HOST") {
		t.Error("expected LEGACY_HOST to be pruned")
	}
	if !contains(removedKeys, "DEBUG_FLAG") {
		t.Error("expected DEBUG_FLAG to be pruned")
	}
	if len(result.Kept) != 5 {
		t.Errorf("expected 5 kept entries, got %d", len(result.Kept))
	}
}

func TestPrune_RemoveCommented(t *testing.T) {
	result := Prune(pruneEntries(), PruneOptions{
		RemoveCommented: true,
	})
	removedKeys := PrunedKeys(result)
	if !contains(removedKeys, "COMMENTED_VAL") {
		t.Error("expected COMMENTED_VAL to be pruned")
	}
	if len(result.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(result.Removed))
	}
}

func TestPrune_NoOptions_KeepsAll(t *testing.T) {
	entries := pruneEntries()
	result := Prune(entries, PruneOptions{})
	if len(result.Kept) != len(entries) {
		t.Errorf("expected all %d entries kept, got %d", len(entries), len(result.Kept))
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(result.Removed))
	}
}

func TestPrunedKeys_ReturnsKeySlice(t *testing.T) {
	result := Prune(pruneEntries(), PruneOptions{
		RemoveKeys: []string{"APP_NAME", "PORT"},
	})
	keys := PrunedKeys(result)
	if len(keys) != 2 {
		t.Fatalf("expected 2 pruned keys, got %d", len(keys))
	}
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
