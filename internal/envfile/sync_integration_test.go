package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func TestSync_RoundTrip(t *testing.T) {
	dir := t.TempDir()

	sourcePath := filepath.Join(dir, ".env.source")
	targetPath := filepath.Join(dir, ".env.target")

	sourceContent := "DB_HOST=localhost\nDB_PASS=secret\nAPP_ENV=production\n"
	targetContent := "DB_HOST=remotehost\nEXISTING=keep\n"

	if err := os.WriteFile(sourcePath, []byte(sourceContent), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte(targetContent), 0o600); err != nil {
		t.Fatal(err)
	}

	source, err := envfile.Parse(sourcePath)
	if err != nil {
		t.Fatalf("parse source: %v", err)
	}
	target, err := envfile.Parse(targetPath)
	if err != nil {
		t.Fatalf("parse target: %v", err)
	}

	result, err := envfile.Sync(source, target, targetPath, envfile.SyncOptions{
		Overwrite: false,
		DryRun:    false,
	})
	if err != nil {
		t.Fatalf("sync: %v", err)
	}

	// DB_HOST already exists in target, should be skipped
	for _, s := range result.Skipped {
		if s == "DB_HOST" {
			goto skippedOK
		}
	}
	t.Error("expected DB_HOST to be skipped")
skippedOK:

	updated, err := envfile.Parse(targetPath)
	if err != nil {
		t.Fatalf("parse updated target: %v", err)
	}

	if updated["DB_HOST"] != "remotehost" {
		t.Errorf("DB_HOST should remain remotehost, got %s", updated["DB_HOST"])
	}
	if updated["EXISTING"] != "keep" {
		t.Errorf("EXISTING should remain, got %s", updated["EXISTING"])
	}
	if updated["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be synced, got %s", updated["APP_ENV"])
	}
}
