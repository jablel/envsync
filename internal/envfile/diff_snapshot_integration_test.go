package envfile

import (
	"path/filepath"
	"testing"
)

func TestDiffSnapshot_SaveLoadDiffRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	// Baseline state
	initial := []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "db.prod.internal"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}

	snap := TakeSnapshot(initial)
	if err := SaveSnapshot(snap, path); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	// Simulate drift: modify one key, add one, remove one
	current := []Entry{
		{Key: "APP_ENV", Value: "staging"},
		{Key: "DB_HOST", Value: "db.prod.internal"},
		{Key: "NEW_FEATURE_FLAG", Value: "true"},
	}

	d, err := DiffSnapshotByPath(current, path)
	if err != nil {
		t.Fatalf("DiffSnapshotByPath: %v", err)
	}

	if len(d.Added) != 1 || d.Added[0].Key != "NEW_FEATURE_FLAG" {
		t.Errorf("expected NEW_FEATURE_FLAG added, got %+v", d.Added)
	}
	if len(d.Removed) != 1 || d.Removed[0].Key != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY removed, got %+v", d.Removed)
	}
	if len(d.Modified) != 1 || d.Modified[0].Old.Value != "production" {
		t.Errorf("expected APP_ENV modified, got %+v", d.Modified)
	}

	// Verify summary formatting
	s := SummarizeDiffSnapshot(d)
	out := FormatDiffSnapshotSummary(s)
	if out == "" {
		t.Error("expected non-empty summary")
	}

	// Apply diff back to initial and verify convergence
	result := ApplySnapshotDiff(initial, d)
	m := make(map[string]string)
	for _, e := range result {
		m[e.Key] = e.Value
	}
	if m["APP_ENV"] != "staging" {
		t.Errorf("APP_ENV should be staging after apply, got %s", m["APP_ENV"])
	}
	if _, ok := m["SECRET_KEY"]; ok {
		t.Error("SECRET_KEY should have been removed")
	}
	if m["NEW_FEATURE_FLAG"] != "true" {
		t.Errorf("NEW_FEATURE_FLAG should be true after apply, got %s", m["NEW_FEATURE_FLAG"])
	}
}
