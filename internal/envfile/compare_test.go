package envfile

import (
	"strings"
	"testing"
)

func baseEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envsync"},
		{Key: "DEBUG", Value: "true"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func otherEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envsync"},
		{Key: "DEBUG", Value: "false"},
		{Key: "API_URL", Value: "https://api.example.com"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestCompare_MatchingKeys(t *testing.T) {
	r := Compare(baseEntries(), otherEntries())
	if len(r.Matching) != 2 {
		t.Errorf("expected 2 matching keys, got %d", len(r.Matching))
	}
}

func TestCompare_MismatchedValues(t *testing.T) {
	r := Compare(baseEntries(), otherEntries())
	vals, ok := r.Mismatched["DEBUG"]
	if !ok {
		t.Fatal("expected DEBUG to be mismatched")
	}
	if vals[0] != "true" || vals[1] != "false" {
		t.Errorf("unexpected mismatch values: %v", vals)
	}
}

func TestCompare_OnlyInBase(t *testing.T) {
	r := Compare(baseEntries(), otherEntries())
	if len(r.OnlyInBase) != 1 || r.OnlyInBase[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST only in base, got %v", r.OnlyInBase)
	}
}

func TestCompare_OnlyInOther(t *testing.T) {
	r := Compare(baseEntries(), otherEntries())
	if len(r.OnlyInOther) != 1 || r.OnlyInOther[0] != "API_URL" {
		t.Errorf("expected API_URL only in other, got %v", r.OnlyInOther)
	}
}

func TestCompare_IsIdentical_False(t *testing.T) {
	r := Compare(baseEntries(), otherEntries())
	if r.IsIdentical() {
		t.Error("expected IsIdentical to return false")
	}
}

func TestCompare_IsIdentical_True(t *testing.T) {
	entries := baseEntries()
	r := Compare(entries, entries)
	if !r.IsIdentical() {
		t.Error("expected IsIdentical to return true for identical inputs")
	}
}

func TestCompare_Summary_ContainsFields(t *testing.T) {
	r := Compare(baseEntries(), otherEntries())
	summary := r.Summary()
	for _, field := range []string{"Matching", "Mismatched", "Only in base", "Only in other"} {
		if !strings.Contains(summary, field) {
			t.Errorf("summary missing field %q", field)
		}
	}
}

func TestCompare_EmptyInputs(t *testing.T) {
	r := Compare(nil, nil)
	if !r.IsIdentical() {
		t.Error("two empty env files should be identical")
	}
}
