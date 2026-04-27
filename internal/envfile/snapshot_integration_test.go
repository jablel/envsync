package envfile

import (
	"path/filepath"
	"testing"
)

func TestSnapshot_SaveLoadDiffRoundTrip(t *testing.T) {
	dir := t.TempDir()

	v1Entries := map[string]string{
		"APP_ENV":  "staging",
		"DB_HOST":  "db.staging.internal",
		"LOG_LEVEL": "debug",
	}
	v2Entries := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "db.prod.internal",
		"LOG_LEVEL": "info",
		"CACHE_TTL": "300",
	}

	v1 := TakeSnapshot("v1", v1Entries)
	v2 := TakeSnapshot("v2", v2Entries)

	v1Path := filepath.Join(dir, "v1.json")
	v2Path := filepath.Join(dir, "v2.json")

	if err := SaveSnapshot(v1Path, v1); err != nil {
		t.Fatalf("save v1: %v", err)
	}
	if err := SaveSnapshot(v2Path, v2); err != nil {
		t.Fatalf("save v2: %v", err)
	}

	loadedV1, err := LoadSnapshot(v1Path)
	if err != nil {
		t.Fatalf("load v1: %v", err)
	}
	loadedV2, err := LoadSnapshot(v2Path)
	if err != nil {
		t.Fatalf("load v2: %v", err)
	}

	result := DiffSnapshots(loadedV1, loadedV2)

	if len(result.Added) != 1 || result.Added[0].Key != "CACHE_TTL" {
		t.Errorf("expected CACHE_TTL added, got %v", result.Added)
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected no removals, got %v", result.Removed)
	}
	if len(result.Modified) != 3 {
		t.Errorf("expected 3 modified keys, got %d: %v", len(result.Modified), result.Modified)
	}
	if !HasSnapshotChanges(result) {
		t.Error("expected changes between v1 and v2")
	}

	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
