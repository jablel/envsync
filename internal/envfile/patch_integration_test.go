package envfile_test

import (
	"os"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func TestPatch_RoundTrip(t *testing.T) {
	original := `APP_ENV=production
DB_HOST=localhost
DB_PORT=5432
SECRET_KEY=supersecret
`
	tmp, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := tmp.WriteString(original); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmp.Close()

	entries, err := envfile.Parse(tmp.Name())
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	instructions := []envfile.PatchInstruction{
		{Op: envfile.PatchSet, Key: "APP_ENV", Value: "staging"},
		{Op: envfile.PatchSet, Key: "NEW_FEATURE_FLAG", Value: "true"},
		{Op: envfile.PatchDelete, Key: "DB_PORT"},
		{Op: envfile.PatchRename, Key: "DB_HOST", NewKey: "DATABASE_HOST"},
	}

	patched, results, err := envfile.Patch(entries, instructions)
	if err != nil {
		t.Fatalf("patch error: %v", err)
	}

	for _, r := range results {
		if !r.Applied {
			t.Errorf("instruction not applied: %+v, reason: %s", r.Instruction, r.Reason)
		}
	}

	keyMap := make(map[string]string)
	for _, e := range patched {
		keyMap[e.Key] = e.Value
	}

	if keyMap["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", keyMap["APP_ENV"])
	}
	if keyMap["NEW_FEATURE_FLAG"] != "true" {
		t.Errorf("expected NEW_FEATURE_FLAG=true, got %q", keyMap["NEW_FEATURE_FLAG"])
	}
	if _, ok := keyMap["DB_PORT"]; ok {
		t.Error("expected DB_PORT to be deleted")
	}
	if _, ok := keyMap["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be renamed away")
	}
	if keyMap["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", keyMap["DATABASE_HOST"])
	}
}
