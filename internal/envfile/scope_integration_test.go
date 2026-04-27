package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func TestScope_ParseAndScopeRoundTrip(t *testing.T) {
	content := `DB_HOST=localhost
DB_PORT=5432
DB_PASS=secret
`
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	entries, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	// Add prefix
	scoped, err := envfile.Scope(entries, envfile.ScopeOptions{Prefix: "STAGING_"})
	if err != nil {
		t.Fatalf("scope add: %v", err)
	}
	if len(scoped.Entries) != 3 {
		t.Fatalf("expected 3 scoped entries, got %d", len(scoped.Entries))
	}

	// Strip prefix back
	stripped, err := envfile.Scope(scoped.Entries, envfile.ScopeOptions{
		Prefix:      "STAGING_",
		Strip:       true,
		SkipNoMatch: true,
	})
	if err != nil {
		t.Fatalf("scope strip: %v", err)
	}
	if len(stripped.Entries) != 3 {
		t.Fatalf("expected 3 stripped entries, got %d", len(stripped.Entries))
	}

	// Keys should match originals
	origMap := make(map[string]string)
	for _, e := range entries {
		origMap[e.Key] = e.Value
	}
	for _, e := range stripped.Entries {
		v, ok := origMap[e.Key]
		if !ok {
			t.Errorf("key %s not found in original", e.Key)
		}
		if v != e.Value {
			t.Errorf("value mismatch for %s: want %s got %s", e.Key, v, e.Value)
		}
	}
}
