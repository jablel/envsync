package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPromote_RoundTrip(t *testing.T) {
	dir := t.TempDir()

	srcContent := "DB_HOST=prod-db.example.com\nDB_PORT=5432\nSECRET_KEY=s3cr3t\n"
	dstContent := "DB_HOST=localhost\nAPP_ENV=staging\n"

	srcFile := filepath.Join(dir, "prod.env")
	dstFile := filepath.Join(dir, "staging.env")

	if err := os.WriteFile(srcFile, []byte(srcContent), 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	if err := os.WriteFile(dstFile, []byte(dstContent), 0o644); err != nil {
		t.Fatalf("write dst: %v", err)
	}

	srcEntries, err := Parse(srcFile)
	if err != nil {
		t.Fatalf("parse src: %v", err)
	}
	dstParsed, err := Parse(dstFile)
	if err != nil {
		t.Fatalf("parse dst: %v", err)
	}

	updated, result, err := Promote(srcEntries, dstParsed, PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("promote: %v", err)
	}

	outFile := filepath.Join(dir, "staging_promoted.env")
	if err := writeEnvFile(outFile, updated); err != nil {
		t.Fatalf("write output: %v", err)
	}

	reloaded, err := Parse(outFile)
	if err != nil {
		t.Fatalf("parse output: %v", err)
	}

	m := entriesToMap(reloaded)
	if m["DB_HOST"] != "prod-db.example.com" {
		t.Errorf("expected overwritten DB_HOST, got %q", m["DB_HOST"])
	}
	if m["APP_ENV"] != "staging" {
		t.Errorf("expected preserved APP_ENV, got %q", m["APP_ENV"])
	}
	if m["SECRET_KEY"] != "s3cr3t" {
		t.Errorf("expected promoted SECRET_KEY, got %q", m["SECRET_KEY"])
	}
	if len(result.Promoted) == 0 && len(result.Overwritten) == 0 {
		t.Error("expected at least some promoted or overwritten keys")
	}
}
