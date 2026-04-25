package envfile

import (
	"testing"
)

func entries(pairs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestMerge_NewKeysAppended(t *testing.T) {
	base := entries("A", "1", "B", "2")
	override := entries("C", "3")
	res, err := Merge(base, override, MergeOptions{Strategy: MergeStrategyKeepBase})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(res.Entries))
	}
	if res.Entries[2].Key != "C" || res.Entries[2].Value != "3" {
		t.Errorf("expected C=3, got %v", res.Entries[2])
	}
}

func TestMerge_KeepBase_OnConflict(t *testing.T) {
	base := entries("A", "original")
	override := entries("A", "changed")
	res, err := Merge(base, override, MergeOptions{Strategy: MergeStrategyKeepBase})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Value != "original" {
		t.Errorf("expected base value 'original', got %q", res.Entries[0].Value)
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "A" {
		t.Errorf("expected conflict on A, got %v", res.Conflicts)
	}
}

func TestMerge_KeepOverride_OnConflict(t *testing.T) {
	base := entries("A", "original")
	override := entries("A", "changed")
	res, err := Merge(base, override, MergeOptions{Strategy: MergeStrategyKeepOverride})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Value != "changed" {
		t.Errorf("expected override value 'changed', got %q", res.Entries[0].Value)
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %v", res.Conflicts)
	}
}

func TestMerge_ErrorStrategy_OnConflict(t *testing.T) {
	base := entries("A", "1")
	override := entries("A", "2")
	_, err := Merge(base, override, MergeOptions{Strategy: MergeStrategyError})
	if err == nil {
		t.Fatal("expected error for conflict, got nil")
	}
}

func TestMerge_SameValue_NoConflict(t *testing.T) {
	base := entries("A", "same")
	override := entries("A", "same")
	res, err := Merge(base, override, MergeOptions{Strategy: MergeStrategyError})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestMerge_SkipBlanks(t *testing.T) {
	base := []Entry{{Key: "A", Value: "1"}, {Key: "", Value: ""}}
	override := []Entry{{Key: "", Value: ""}, {Key: "B", Value: "2"}}
	res, err := Merge(base, override, MergeOptions{Strategy: MergeStrategyKeepBase, SkipBlanks: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries after skipping blanks, got %d", len(res.Entries))
	}
}
