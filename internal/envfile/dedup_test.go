package envfile

import (
	"testing"
)

func TestDedup_KeepFirst(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
		{Key: "HOST", Value: "remotehost"},
	}
	res, err := Dedup(entries, DedupKeepFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Value != "localhost" {
		t.Errorf("expected localhost, got %s", res.Entries[0].Value)
	}
	if res.Duplicates["HOST"] != 1 {
		t.Errorf("expected 1 duplicate for HOST, got %d", res.Duplicates["HOST"])
	}
}

func TestDedup_KeepLast(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
		{Key: "HOST", Value: "remotehost"},
	}
	res, err := Dedup(entries, DedupKeepLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Value != "remotehost" {
		t.Errorf("expected remotehost, got %s", res.Entries[0].Value)
	}
	if len(res.Removed) != 1 || res.Removed[0].Value != "localhost" {
		t.Errorf("expected removed entry with value localhost")
	}
}

func TestDedup_ErrorStrategy(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "HOST", Value: "remotehost"},
	}
	_, err := Dedup(entries, DedupError)
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestDedup_NoDuplicates(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
	}
	res, err := Dedup(entries, DedupKeepFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removed entries")
	}
}

func TestHasDuplicates_True(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "A", Value: "2"},
	}
	if !HasDuplicates(entries) {
		t.Error("expected HasDuplicates to return true")
	}
}

func TestHasDuplicates_False(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	if HasDuplicates(entries) {
		t.Error("expected HasDuplicates to return false")
	}
}
