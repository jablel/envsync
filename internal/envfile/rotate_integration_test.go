package envfile

import (
	"fmt"
	"strings"
	"testing"
)

func TestRotate_RoundTrip(t *testing.T) {
	// Write an env file, parse it, rotate secrets, verify values changed.
	content := "APP_NAME=myapp\nDB_PASSWORD=hunter2\nAPI_SECRET=topsecret\nPORT=3000\n"
	f := writeTempEnv(t, content)

	entries, err := Parse(f)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	masker := NewMasker(nil)
	counter := 0
	gen := func(key string) (string, error) {
		counter++
		return fmt.Sprintf("rotated-value-%d", counter), nil
	}

	updated, result, err := Rotate(entries, masker, RotateOptions{Generator: gen})
	if err != nil {
		t.Fatalf("rotate error: %v", err)
	}

	if len(result.Rotated) == 0 {
		t.Fatal("expected at least one key to be rotated")
	}

	updatedMap := make(map[string]string)
	for _, e := range updated {
		updatedMap[e.Key] = e.Value
	}

	for _, k := range result.Rotated {
		if !strings.HasPrefix(updatedMap[k], "rotated-value-") {
			t.Errorf("key %s not updated after rotation, value=%s", k, updatedMap[k])
		}
	}

	// Non-sensitive keys must remain unchanged.
	if updatedMap["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should be unchanged, got %s", updatedMap["APP_NAME"])
	}
	if updatedMap["PORT"] != "3000" {
		t.Errorf("PORT should be unchanged, got %s", updatedMap["PORT"])
	}

	// Verify summary and instructions are non-empty.
	summary := RotateSummary(result)
	if !strings.Contains(summary, "Rotated:") {
		t.Errorf("summary missing Rotated section: %s", summary)
	}

	instructions := FormatRotateInstructions(entries, updated, masker)
	if instructions == "" {
		t.Error("expected non-empty instructions after rotation")
	}
}
