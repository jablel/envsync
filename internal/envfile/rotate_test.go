package envfile

import (
	"fmt"
	"strings"
	"testing"
)

func rotateEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "old-pass"},
		{Key: "API_SECRET", Value: "old-secret"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestRotate_AllSensitiveKeys(t *testing.T) {
	masker := NewMasker(nil)
	gen := func(key string) (string, error) { return "new-" + strings.ToLower(key), nil }

	updated, result, err := Rotate(rotateEntries(), masker, RotateOptions{Generator: gen})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Rotated) == 0 {
		t.Fatal("expected at least one rotated key")
	}
	for _, e := range updated {
		if masker.IsSensitive(e.Key) && !strings.HasPrefix(e.Value, "new-") {
			t.Errorf("key %s was not rotated, value=%s", e.Key, e.Value)
		}
	}
}

func TestRotate_SpecificKeys(t *testing.T) {
	masker := NewMasker(nil)
	gen := func(key string) (string, error) { return "rotated", nil }

	updated, result, err := Rotate(rotateEntries(), masker, RotateOptions{
		Keys:      []string{"DB_PASSWORD"},
		Generator: gen,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Rotated) != 1 || result.Rotated[0] != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD rotated, got %v", result.Rotated)
	}
	for _, e := range updated {
		if e.Key == "API_SECRET" && e.Value != "old-secret" {
			t.Error("API_SECRET should not have been rotated")
		}
		if e.Key == "DB_PASSWORD" && e.Value != "rotated" {
			t.Errorf("DB_PASSWORD not rotated, got %s", e.Value)
		}
	}
}

func TestRotate_DryRun(t *testing.T) {
	masker := NewMasker(nil)
	gen := func(key string) (string, error) { return "should-not-apply", nil }

	updated, result, err := Rotate(rotateEntries(), masker, RotateOptions{
		Generator: gen,
		DryRun:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Rotated) == 0 {
		t.Fatal("dry run should still report rotated keys")
	}
	for _, e := range updated {
		if e.Value == "should-not-apply" {
			t.Errorf("dry run must not modify values, key=%s", e.Key)
		}
	}
}

func TestRotate_GeneratorError(t *testing.T) {
	masker := NewMasker(nil)
	gen := func(key string) (string, error) { return "", fmt.Errorf("vault unavailable") }

	_, result, err := Rotate(rotateEntries(), masker, RotateOptions{
		Keys:      []string{"DB_PASSWORD"},
		Generator: gen,
	})
	if err == nil {
		t.Fatal("expected error from generator failure")
	}
	if _, ok := result.Errors["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD in error map")
	}
}

func TestRotate_NilGenerator(t *testing.T) {
	masker := NewMasker(nil)
	_, _, err := Rotate(rotateEntries(), masker, RotateOptions{})
	if err == nil {
		t.Fatal("expected error when generator is nil")
	}
}
