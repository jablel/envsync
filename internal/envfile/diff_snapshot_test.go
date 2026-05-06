package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiffSnapshot_Added(t *testing.T) {
	snap := Snapshot{Entries: []Entry{{Key: "A", Value: "1"}}}
	current := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	d := DiffSnapshot(current, snap)
	if len(d.Added) != 1 || d.Added[0].Key != "B" {
		t.Errorf("expected B to be added, got %+v", d.Added)
	}
}

func TestDiffSnapshot_Removed(t *testing.T) {
	snap := Snapshot{Entries: []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}}
	current := []Entry{{Key: "A", Value: "1"}}
	d := DiffSnapshot(current, snap)
	if len(d.Removed) != 1 || d.Removed[0].Key != "B" {
		t.Errorf("expected B to be removed, got %+v", d.Removed)
	}
}

func TestDiffSnapshot_Modified(t *testing.T) {
	snap := Snapshot{Entries: []Entry{{Key: "A", Value: "old"}}}
	current := []Entry{{Key: "A", Value: "new"}}
	d := DiffSnapshot(current, snap)
	if len(d.Modified) != 1 || d.Modified[0].New.Value != "new" {
		t.Errorf("expected A to be modified, got %+v", d.Modified)
	}
}

func TestDiffSnapshotByPath_LoadsAndDiffs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	snap := Snapshot{Entries: []Entry{{Key: "X", Value: "10"}}}
	if err := SaveSnapshot(snap, path); err != nil {
		t.Fatalf("save: %v", err)
	}
	current := []Entry{{Key: "X", Value: "10"}, {Key: "Y", Value: "20"}}
	d, err := DiffSnapshotByPath(current, path)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if len(d.Added) != 1 || d.Added[0].Key != "Y" {
		t.Errorf("expected Y added, got %+v", d.Added)
	}
}

func TestDiffSnapshotByPath_MissingFile(t *testing.T) {
	_, err := DiffSnapshotByPath(nil, "/nonexistent/snap.json")
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}

func TestApplySnapshotDiff_AddsAndRemoves(t *testing.T) {
	base := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	snap := Snapshot{Entries: []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}, {Key: "C", Value: "3"}}}
	current := []Entry{{Key: "A", Value: "99"}, {Key: "D", Value: "4"}}
	d := DiffSnapshot(current, snap)
	result := ApplySnapshotDiff(base, d)
	m := make(map[string]string)
	for _, e := range result {
		m[e.Key] = e.Value
	}
	if m["A"] != "99" {
		t.Errorf("expected A=99, got %s", m["A"])
	}
	if _, ok := m["C"]; ok {
		t.Error("expected C to be removed")
	}
	if m["D"] != "4" {
		t.Errorf("expected D=4, got %s", m["D"])
	}
}
