package envfile

import "testing"

func baseSnap() Snapshot {
	return TakeSnapshot("base", map[string]string{
		"APP_ENV":  "development",
		"DB_HOST":  "localhost",
		"API_KEY":  "secret123",
	})
}

func TestDiffSnapshots_Added(t *testing.T) {
	from := baseSnap()
	to := TakeSnapshot("next", map[string]string{
		"APP_ENV": "development",
		"DB_HOST": "localhost",
		"API_KEY": "secret123",
		"NEW_KEY": "newval",
	})
	result := DiffSnapshots(from, to)
	if len(result.Added) != 1 || result.Added[0].Key != "NEW_KEY" {
		t.Errorf("expected 1 added key NEW_KEY, got %v", result.Added)
	}
}

func TestDiffSnapshots_Removed(t *testing.T) {
	from := baseSnap()
	to := TakeSnapshot("next", map[string]string{
		"APP_ENV": "development",
		"DB_HOST": "localhost",
	})
	result := DiffSnapshots(from, to)
	if len(result.Removed) != 1 || result.Removed[0].Key != "API_KEY" {
		t.Errorf("expected 1 removed key API_KEY, got %v", result.Removed)
	}
}

func TestDiffSnapshots_Modified(t *testing.T) {
	from := baseSnap()
	to := TakeSnapshot("next", map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "localhost",
		"API_KEY": "secret123",
	})
	result := DiffSnapshots(from, to)
	if len(result.Modified) != 1 || result.Modified[0].Key != "APP_ENV" {
		t.Errorf("expected 1 modified key APP_ENV, got %v", result.Modified)
	}
	if result.Modified[0].OldValue != "development" || result.Modified[0].NewValue != "production" {
		t.Errorf("unexpected modified values: %+v", result.Modified[0])
	}
}

func TestDiffSnapshots_Unchanged(t *testing.T) {
	from := baseSnap()
	to := baseSnap()
	result := DiffSnapshots(from, to)
	if len(result.Unchanged) != 3 {
		t.Errorf("expected 3 unchanged, got %d", len(result.Unchanged))
	}
	if HasSnapshotChanges(result) {
		t.Error("expected no changes")
	}
}

func TestDiffSnapshots_Summary(t *testing.T) {
	from := baseSnap()
	to := TakeSnapshot("prod", map[string]string{"APP_ENV": "production", "NEW": "val"})
	result := DiffSnapshots(from, to)
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestHasSnapshotChanges_True(t *testing.T) {
	from := baseSnap()
	to := TakeSnapshot("other", map[string]string{"ONLY_THIS": "val"})
	if !HasSnapshotChanges(DiffSnapshots(from, to)) {
		t.Error("expected changes to be detected")
	}
}
