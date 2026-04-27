package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTakeSnapshot_Basic(t *testing.T) {
	entries := map[string]string{
		"APP_NAME": "envsync",
		"PORT":     "8080",
	}
	snap := TakeSnapshot("test", entries)
	if snap.Label != "test" {
		t.Errorf("expected label 'test', got %q", snap.Label)
	}
	if len(snap.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap.Entries))
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestTakeSnapshot_SortedKeys(t *testing.T) {
	entries := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	snap := TakeSnapshot("sorted", entries)
	keys := []string{snap.Entries[0].Key, snap.Entries[1].Key, snap.Entries[2].Key}
	if keys[0] != "A_KEY" || keys[1] != "M_KEY" || keys[2] != "Z_KEY" {
		t.Errorf("entries not sorted: %v", keys)
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	original := TakeSnapshot("prod", map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"})
	original.Timestamp = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	if err := SaveSnapshot(path, original); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}
	if loaded.Label != "prod" {
		t.Errorf("label mismatch: got %q", loaded.Label)
	}
	if len(loaded.Entries) != 2 {
		t.Errorf("entry count mismatch: got %d", len(loaded.Entries))
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSnapshot_ToMap(t *testing.T) {
	snap := TakeSnapshot("dev", map[string]string{"FOO": "bar", "BAZ": "qux"})
	m := snap.ToMap()
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Errorf("ToMap result mismatch: %v", m)
	}
}

func TestSaveSnapshot_InvalidPath(t *testing.T) {
	snap := TakeSnapshot("x", map[string]string{})
	err := SaveSnapshot("/no/such/dir/snap.json", snap)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestSnapshot_EmptyEntries(t *testing.T) {
	snap := TakeSnapshot("empty", map[string]string{})
	if len(snap.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(snap.Entries))
	}
	_ = os.Getenv("CI") // no-op to avoid unused import
}
