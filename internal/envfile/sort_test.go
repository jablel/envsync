package envfile

import (
	"testing"
)

func sortTestEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "AWS_SECRET", Value: "topsecret"},
		{Key: "APP_NAME", Value: "envsync"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "AWS_KEY", Value: "keyval"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestSort_Alpha(t *testing.T) {
	result := Sort(sortTestEntries(), DefaultSortOptions())
	expected := []string{"APP_NAME", "AWS_KEY", "AWS_SECRET", "DB_HOST", "DB_PORT", "PORT"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("pos %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestSort_AlphaDesc(t *testing.T) {
	result := Sort(sortTestEntries(), SortOptions{Order: SortAlphaDesc})
	expected := []string{"PORT", "DB_PORT", "DB_HOST", "AWS_SECRET", "AWS_KEY", "APP_NAME"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("pos %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestSort_ByPrefix(t *testing.T) {
	result := Sort(sortTestEntries(), SortOptions{Order: SortByPrefix})
	// Groups: APP, AWS, DB, PORT
	if result[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME first, got %q", result[0].Key)
	}
	if result[1].Key != "AWS_KEY" || result[2].Key != "AWS_SECRET" {
		t.Errorf("expected AWS group sorted, got %q %q", result[1].Key, result[2].Key)
	}
	if result[3].Key != "DB_HOST" || result[4].Key != "DB_PORT" {
		t.Errorf("expected DB group sorted, got %q %q", result[3].Key, result[4].Key)
	}
	if result[5].Key != "PORT" {
		t.Errorf("expected PORT last, got %q", result[5].Key)
	}
}

func TestSort_SecretLast(t *testing.T) {
	masker := NewMasker(nil) // AWS_SECRET and AWS_KEY should be sensitive
	result := Sort(sortTestEntries(), SortOptions{
		Order:  SortSecretLast,
		Masker: masker,
	})
	// sensitive keys should appear after non-sensitive
	foundSensitive := false
	for _, e := range result {
		isSensitive := masker.IsSensitive(e.Key)
		if foundSensitive && !isSensitive {
			t.Errorf("non-sensitive key %q appeared after sensitive key", e.Key)
		}
		if isSensitive {
			foundSensitive = true
		}
	}
}

func TestSort_DoesNotMutateInput(t *testing.T) {
	input := sortTestEntries()
	origFirst := input[0].Key
	Sort(input, DefaultSortOptions())
	if input[0].Key != origFirst {
		t.Errorf("input was mutated: first key changed from %q to %q", origFirst, input[0].Key)
	}
}

func TestSortedKeys_ReturnsKeysOnly(t *testing.T) {
	keys := SortedKeys(sortTestEntries(), DefaultSortOptions())
	expected := []string{"APP_NAME", "AWS_KEY", "AWS_SECRET", "DB_HOST", "DB_PORT", "PORT"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, expected[i])
		}
	}
}
