package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilter_ParseAndFilterRoundTrip(t *testing.T) {
	content := `APP_NAME=myapp
APP_ENV=staging
DB_HOST=db.internal
DB_PASSWORD=supersecret
API_KEY=key123
DEBUG=true
`
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	entries, err := Parse(p)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	// Filter to only APP_ prefixed, non-secret entries.
	m := NewMasker(nil)
	result, err := Filter(entries, FilterOptions{
		Prefix:         "APP_",
		ExcludeSecrets: true,
		Masker:         m,
	})
	if err != nil {
		t.Fatalf("filter: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 APP_ entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key != "APP_NAME" && e.Key != "APP_ENV" {
			t.Errorf("unexpected key in result: %q", e.Key)
		}
	}

	// Verify DB_ prefix filtering.
	dbResult, err := Filter(entries, FilterOptions{Prefix: "DB_"})
	if err != nil {
		t.Fatalf("filter db: %v", err)
	}
	if len(dbResult) != 2 {
		t.Fatalf("expected 2 DB_ entries, got %d", len(dbResult))
	}
}
