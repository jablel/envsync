package envfile

import (
	"testing"
)

var groupEntries = []Entry{
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "DB_PORT", Value: "5432"},
	{Key: "AWS_ACCESS_KEY", Value: "AKIA..."},
	{Key: "AWS_SECRET", Value: "secret"},
	{Key: "APP_ENV", Value: "production"},
	{Key: "LOG_LEVEL", Value: "info"},
}

func TestGroup_ByPrefix(t *testing.T) {
	opts := GroupOptions{
		Prefixes:     []string{"DB", "AWS"},
		DefaultGroup: "other",
	}
	g := Group(groupEntries, opts)

	if len(g["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(g["DB"]))
	}
	if len(g["AWS"]) != 2 {
		t.Errorf("expected 2 AWS entries, got %d", len(g["AWS"]))
	}
	if len(g["other"]) != 2 {
		t.Errorf("expected 2 other entries, got %d", len(g["other"]))
	}
}

func TestGroup_DefaultGroupOnly(t *testing.T) {
	opts := DefaultGroupOptions()
	g := Group(groupEntries, opts)

	if len(g) != 1 {
		t.Errorf("expected 1 group, got %d", len(g))
	}
	if len(g["general"]) != len(groupEntries) {
		t.Errorf("expected all entries in general group")
	}
}

func TestGroup_EmptyEntries(t *testing.T) {
	opts := GroupOptions{Prefixes: []string{"DB"}, DefaultGroup: "other"}
	g := Group([]Entry{}, opts)
	if len(g) != 0 {
		t.Errorf("expected empty result, got %d groups", len(g))
	}
}

func TestGroupNames_Sorted(t *testing.T) {
	g := GroupedEntries{
		"zebra": {},
		"apple": {},
		"mango": {},
	}
	names := GroupNames(g)
	expected := []string{"apple", "mango", "zebra"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, n)
		}
	}
}

func TestFlatten_PreservesAllEntries(t *testing.T) {
	opts := GroupOptions{
		Prefixes:     []string{"DB", "AWS"},
		DefaultGroup: "other",
	}
	g := Group(groupEntries, opts)
	flat := Flatten(g)

	if len(flat) != len(groupEntries) {
		t.Errorf("expected %d entries after flatten, got %d", len(groupEntries), len(flat))
	}
}

func TestGroup_FirstPrefixWins(t *testing.T) {
	entries := []Entry{
		{Key: "DB_AWS_HOST", Value: "x"},
	}
	opts := GroupOptions{
		Prefixes:     []string{"DB", "AWS"},
		DefaultGroup: "other",
	}
	g := Group(entries, opts)
	if len(g["DB"]) != 1 {
		t.Errorf("expected DB group to win for DB_AWS_HOST")
	}
	if len(g["AWS"]) != 0 {
		t.Errorf("expected AWS group to be empty")
	}
}
