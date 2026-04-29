package envfile

import (
	"testing"
)

var filterEntries = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "APP_ENV", Value: "production"},
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "DB_PASSWORD", Value: "secret123"},
	{Key: "API_KEY", Value: "abc"},
	{Key: "DEBUG", Value: "false"},
}

func TestFilter_ByExplicitKeys(t *testing.T) {
	result, err := Filter(filterEntries, FilterOptions{Keys: []string{"APP_NAME", "DEBUG"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Key != "APP_NAME" || result[1].Key != "DEBUG" {
		t.Errorf("unexpected keys: %v", result)
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	result, err := Filter(filterEntries, FilterOptions{Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFilter_ByPattern(t *testing.T) {
	result, err := Filter(filterEntries, FilterOptions{Pattern: "^APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	_, err := Filter(filterEntries, FilterOptions{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestFilter_ExcludeSecrets(t *testing.T) {
	m := NewMasker(nil)
	result, err := Filter(filterEntries, FilterOptions{
		ExcludeSecrets: true,
		Masker:         m,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range result {
		if m.IsSensitive(e.Key) {
			t.Errorf("sensitive key %q should have been excluded", e.Key)
		}
	}
}

func TestFilter_CombinedPrefixAndPattern(t *testing.T) {
	result, err := Filter(filterEntries, FilterOptions{
		Prefix:  "APP_",
		Pattern: "ENV$",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].Key != "APP_ENV" {
		t.Errorf("expected only APP_ENV, got %v", result)
	}
}

func TestFilterKeys_ReturnsOnlyKeys(t *testing.T) {
	keys, err := FilterKeys(filterEntries, FilterOptions{Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "DB_HOST" || keys[1] != "DB_PASSWORD" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestFilter_EmptyOptions_ReturnsAll(t *testing.T) {
	result, err := Filter(filterEntries, FilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(filterEntries) {
		t.Errorf("expected %d entries, got %d", len(filterEntries), len(result))
	}
}
