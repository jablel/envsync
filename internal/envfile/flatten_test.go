package envfile

import (
	"testing"
)

func TestFlattenKeys_PassthroughSimpleKeys(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
	}
	opts := DefaultFlattenOptions()
	result, err := FlattenKeys(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Key != "HOST" || result[1].Key != "PORT" {
		t.Errorf("unexpected keys: %v", result)
	}
}

func TestFlattenKeys_NestedKeysPreserved(t *testing.T) {
	entries := []Entry{
		{Key: "DB__HOST", Value: "localhost"},
		{Key: "DB__PORT", Value: "5432"},
		{Key: "APP__NAME", Value: "envsync"},
	}
	opts := DefaultFlattenOptions()
	result, err := FlattenKeys(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestFlattenKeys_PrefixFilter(t *testing.T) {
	entries := []Entry{
		{Key: "DB__HOST", Value: "localhost"},
		{Key: "APP__NAME", Value: "envsync"},
		{Key: "DEBUG", Value: "true"},
	}
	opts := FlattenOptions{Separator: "__", Prefix: "DB"}
	result, err := FlattenKeys(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestFlattenKeys_CustomSeparator(t *testing.T) {
	entries := []Entry{
		{Key: "DB.HOST", Value: "localhost"},
		{Key: "DB.PORT", Value: "5432"},
	}
	opts := FlattenOptions{Separator: "."}
	result, err := FlattenKeys(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestExpandKeys_ValidEntries(t *testing.T) {
	entries := []Entry{
		{Key: "DB__HOST", Value: "localhost"},
		{Key: "SIMPLE", Value: "val"},
	}
	opts := DefaultFlattenOptions()
	result, err := ExpandKeys(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestExpandKeys_EmptySegmentReturnsError(t *testing.T) {
	entries := []Entry{
		{Key: "DB____HOST", Value: "localhost"},
	}
	opts := DefaultFlattenOptions()
	_, err := ExpandKeys(entries, opts)
	if err == nil {
		t.Fatal("expected error for empty segment, got nil")
	}
}

func TestFlattenKeys_EmptyInput(t *testing.T) {
	opts := DefaultFlattenOptions()
	result, err := FlattenKeys(nil, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
