package envfile

import (
	"testing"
)

func TestDiff_AddedKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	entries := Diff(base, target)
	found := findEntry(entries, "NEW_KEY")
	if found == nil {
		t.Fatal("expected NEW_KEY in diff entries")
	}
	if found.Status != StatusAdded {
		t.Errorf("expected status added, got %s", found.Status)
	}
	if found.NewValue != "value" {
		t.Errorf("expected NewValue 'value', got %s", found.NewValue)
	}
}

func TestDiff_RemovedKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD_KEY": "old"}
	target := map[string]string{"FOO": "bar"}

	entries := Diff(base, target)
	found := findEntry(entries, "OLD_KEY")
	if found == nil {
		t.Fatal("expected OLD_KEY in diff entries")
	}
	if found.Status != StatusRemoved {
		t.Errorf("expected status removed, got %s", found.Status)
	}
	if found.OldValue != "old" {
		t.Errorf("expected OldValue 'old', got %s", found.OldValue)
	}
}

func TestDiff_ModifiedKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "baz"}

	entries := Diff(base, target)
	found := findEntry(entries, "FOO")
	if found == nil {
		t.Fatal("expected FOO in diff entries")
	}
	if found.Status != StatusModified {
		t.Errorf("expected status modified, got %s", found.Status)
	}
	if found.OldValue != "bar" || found.NewValue != "baz" {
		t.Errorf("unexpected values: old=%s new=%s", found.OldValue, found.NewValue)
	}
}

func TestDiff_UnchangedKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar"}

	entries := Diff(base, target)
	found := findEntry(entries, "FOO")
	if found == nil {
		t.Fatal("expected FOO in diff entries")
	}
	if found.Status != StatusUnchanged {
		t.Errorf("expected status unchanged, got %s", found.Status)
	}
}

func TestHasChanges(t *testing.T) {
	noChange := []DiffEntry{{Key: "A", Status: StatusUnchanged}}
	if HasChanges(noChange) {
		t.Error("expected no changes")
	}

	withChange := []DiffEntry{{Key: "A", Status: StatusAdded}}
	if !HasChanges(withChange) {
		t.Error("expected changes to be detected")
	}
}

func findEntry(entries []DiffEntry, key string) *DiffEntry {
	for i := range entries {
		if entries[i].Key == key {
			return &entries[i]
		}
	}
	return nil
}
