package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGroup_ParseAndGroupRoundTrip(t *testing.T) {
	content := `DB_HOST=localhost
DB_PORT=5432
AWS_ACCESS_KEY=AKIA123
AWS_REGION=us-east-1
APP_ENV=staging
LOG_LEVEL=debug
`
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}

	entries, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	opts := GroupOptions{
		Prefixes:     []string{"DB", "AWS"},
		DefaultGroup: "app",
	}
	g := Group(entries, opts)

	if len(g["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(g["DB"]))
	}
	if len(g["AWS"]) != 2 {
		t.Errorf("expected 2 AWS entries, got %d", len(g["AWS"]))
	}
	if len(g["app"]) != 2 {
		t.Errorf("expected 2 app entries, got %d", len(g["app"]))
	}

	// Flatten and verify total count is preserved
	flat := Flatten(g)
	if len(flat) != len(entries) {
		t.Errorf("flatten lost entries: got %d, want %d", len(flat), len(entries))
	}

	// Verify all original keys are present after round-trip
	keySet := make(map[string]bool)
	for _, e := range flat {
		keySet[e.Key] = true
	}
	for _, e := range entries {
		if !keySet[e.Key] {
			t.Errorf("key %q lost during group/flatten round-trip", e.Key)
		}
	}
}
