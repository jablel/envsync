package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTemplate_RoundTrip generates a template from a parsed .env file,
// writes it to disk, re-parses it, and applies values to verify round-trip.
func TestTemplate_RoundTrip(t *testing.T) {
	src := "APP_ENV=production\nSECRET_KEY=abc123\nLOG_LEVEL=debug\n"
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(srcPath, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	entries, err := Parse(srcPath)
	if err != nil {
		t.Fatal(err)
	}

	masker := NewMasker(nil)
	tplLines := GenerateTemplate(entries, masker)

	tplPath := filepath.Join(tmpDir, ".env.template")
	if err := os.WriteFile(tplPath, []byte(strings.Join(tplLines, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	tplEntries, err := Parse(tplPath)
	if err != nil {
		t.Fatal(err)
	}

	// SECRET_KEY should be blank in template
	blankKeys := TemplateKeys(tplEntries)
	if len(blankKeys) == 0 {
		t.Fatal("expected at least one blank key in template")
	}

	values := map[string]string{"SECRET_KEY": "newsecret"}
	result, errs := ApplyTemplate(tplEntries, values)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors applying template: %v", errs)
	}

	resultMap := make(map[string]string)
	for _, e := range result {
		resultMap[e.Key] = e.Value
	}

	if resultMap["SECRET_KEY"] != "newsecret" {
		t.Errorf("expected newsecret, got %s", resultMap["SECRET_KEY"])
	}
	if resultMap["LOG_LEVEL"] != "debug" {
		t.Errorf("expected debug, got %s", resultMap["LOG_LEVEL"])
	}
}
