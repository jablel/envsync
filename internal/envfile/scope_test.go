package envfile

import (
	"testing"
)

func scopeEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "envsync"},
	}
}

func TestScope_AddPrefix(t *testing.T) {
	result, err := Scope(scopeEntries(), ScopeOptions{Prefix: "PROD_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result.Entries))
	}
	if result.Entries[0].Key != "PROD_DB_HOST" {
		t.Errorf("expected PROD_DB_HOST, got %s", result.Entries[0].Key)
	}
	if result.Entries[2].Key != "PROD_APP_NAME" {
		t.Errorf("expected PROD_APP_NAME, got %s", result.Entries[2].Key)
	}
}

func TestScope_StripPrefix_Matching(t *testing.T) {
	result, err := Scope(scopeEntries(), ScopeOptions{Prefix: "DB_", Strip: true, SkipNoMatch: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
	if result.Entries[0].Key != "HOST" {
		t.Errorf("expected HOST, got %s", result.Entries[0].Key)
	}
	if len(result.Skipped) != 1 || result.Skipped[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME in skipped")
	}
}

func TestScope_StripPrefix_NoSkip(t *testing.T) {
	result, err := Scope(scopeEntries(), ScopeOptions{Prefix: "DB_", Strip: true, SkipNoMatch: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// All 3 entries returned: 2 stripped + 1 unchanged
	if len(result.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result.Entries))
	}
	if result.Entries[2].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME unchanged, got %s", result.Entries[2].Key)
	}
}

func TestScope_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := Scope(scopeEntries(), ScopeOptions{Prefix: ""})
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestScope_PreservesValues(t *testing.T) {
	result, err := Scope(scopeEntries(), ScopeOptions{Prefix: "X_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, e := range result.Entries {
		if e.Value != scopeEntries()[i].Value {
			t.Errorf("value mismatch at %d: got %s", i, e.Value)
		}
	}
}

func TestScopeKeys_ReturnsKeys(t *testing.T) {
	keys, err := ScopeKeys(scopeEntries(), ScopeOptions{Prefix: "ENV_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"ENV_DB_HOST", "ENV_DB_PORT", "ENV_APP_NAME"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("key[%d]: expected %s, got %s", i, expected[i], k)
		}
	}
}
