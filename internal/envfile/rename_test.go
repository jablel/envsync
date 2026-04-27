package envfile

import (
	"testing"
)

func baseRenameEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_URL", Value: "postgres://localhost/dev"},
	}
}

func TestRename_Success(t *testing.T) {
	entries, result, err := Rename(baseRenameEntries(), "APP_PORT", "SERVER_PORT", RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Renamed {
		t.Fatalf("expected Renamed=true, got false (reason: %s)", result.Reason)
	}
	for _, e := range entries {
		if e.Key == "APP_PORT" {
			t.Errorf("old key APP_PORT still present after rename")
		}
	}
	found := false
	for _, e := range entries {
		if e.Key == "SERVER_PORT" && e.Value == "8080" {
			found = true
		}
	}
	if !found {
		t.Errorf("new key SERVER_PORT with value 8080 not found")
	}
}

func TestRename_KeyNotFound(t *testing.T) {
	_, result, err := Rename(baseRenameEntries(), "MISSING", "NEW_KEY", RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Renamed {
		t.Errorf("expected Renamed=false when key not found")
	}
	if result.Reason != "key not found" {
		t.Errorf("unexpected reason: %s", result.Reason)
	}
}

func TestRename_IdenticalKeys(t *testing.T) {
	_, result, err := Rename(baseRenameEntries(), "APP_HOST", "APP_HOST", RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Renamed {
		t.Errorf("expected Renamed=false for identical keys")
	}
}

func TestRename_ConflictWithoutOverwrite(t *testing.T) {
	_, _, err := Rename(baseRenameEntries(), "APP_PORT", "DB_URL", RenameOptions{Overwrite: false})
	if err == nil {
		t.Fatal("expected error when target key exists and Overwrite=false")
	}
}

func TestRename_ConflictWithOverwrite(t *testing.T) {
	entries, result, err := Rename(baseRenameEntries(), "APP_PORT", "DB_URL", RenameOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Renamed {
		t.Fatalf("expected Renamed=true")
	}
	count := 0
	for _, e := range entries {
		if e.Key == "DB_URL" {
			count++
			if e.Value != "8080" {
				t.Errorf("expected overwritten value 8080, got %s", e.Value)
			}
		}
	}
	if count != 1 {
		t.Errorf("expected exactly one DB_URL entry, got %d", count)
	}
}

func TestRename_EmptyKeys(t *testing.T) {
	if _, _, err := Rename(baseRenameEntries(), "", "NEW", RenameOptions{}); err == nil {
		t.Error("expected error for empty oldKey")
	}
	if _, _, err := Rename(baseRenameEntries(), "APP_HOST", "", RenameOptions{}); err == nil {
		t.Error("expected error for empty newKey")
	}
}
