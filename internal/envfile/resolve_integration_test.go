package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolve_ParseAndResolveRoundTrip(t *testing.T) {
	baseContent := "APP_NAME=myapp\nDB_HOST=localhost\nDB_PORT=5432\n"
	overrideContent := "DB_HOST=prod-db.example.com\nNEW_KEY=hello\n"

	dir := t.TempDir()
	baseFile := filepath.Join(dir, "base.env")
	overrideFile := filepath.Join(dir, "override.env")

	if err := os.WriteFile(baseFile, []byte(baseContent), 0o600); err != nil {
		t.Fatalf("write base: %v", err)
	}
	if err := os.WriteFile(overrideFile, []byte(overrideContent), 0o600); err != nil {
		t.Fatalf("write override: %v", err)
	}

	baseEntries, err := Parse(baseFile)
	if err != nil {
		t.Fatalf("parse base: %v", err)
	}
	overrideEntries, err := Parse(overrideFile)
	if err != nil {
		t.Fatalf("parse override: %v", err)
	}

	opts := DefaultResolveOptions()
	r, err := Resolve(baseEntries, overrideEntries, opts)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}

	if r.Applied != 2 {
		t.Errorf("expected 2 applied, got %d", r.Applied)
	}

	m := toMap(r.Entries)
	if m["DB_HOST"] != "prod-db.example.com" {
		t.Errorf("expected overridden DB_HOST, got %q", m["DB_HOST"])
	}
	if m["APP_NAME"] != "myapp" {
		t.Errorf("expected unchanged APP_NAME, got %q", m["APP_NAME"])
	}
	if m["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY, got %q", m["NEW_KEY"])
	}

	summary := ResolveSummary(r)
	if summary == "" {
		t.Error("expected non-empty summary")
	}

	formatted := FormatResolvedEntries(r, baseEntries)
	if formatted == "" {
		t.Error("expected non-empty formatted output")
	}
}
