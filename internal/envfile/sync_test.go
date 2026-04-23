package envfile

import (
	"os"
	"testing"
)

func TestSync_AppliesNewKeys(t *testing.T) {
	src := map[string]string{"NEW_KEY": "value1", "ANOTHER": "value2"}
	dst := map[string]string{}

	tmpFile, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	result, err := Sync(src, dst, tmpFile.Name(), SyncOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d", len(result.Applied))
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(result.Skipped))
	}
}

func TestSync_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"KEY": "new"}
	dst := map[string]string{"KEY": "old"}

	tmpFile, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	result, err := Sync(src, dst, tmpFile.Name(), SyncOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
}

func TestSync_OverwriteExistingKeys(t *testing.T) {
	src := map[string]string{"KEY": "new"}
	dst := map[string]string{"KEY": "old"}

	tmpFile, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	result, err := Sync(src, dst, tmpFile.Name(), SyncOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(result.Applied))
	}

	parsed, err := Parse(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if parsed["KEY"] != "new" {
		t.Errorf("expected KEY=new, got %s", parsed["KEY"])
	}
}

func TestSync_DryRunDoesNotWrite(t *testing.T) {
	src := map[string]string{"DRY_KEY": "val"}
	dst := map[string]string{}

	tmpFile, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	_, err = Sync(src, dst, tmpFile.Name(), SyncOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsed, err := Parse(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := parsed["DRY_KEY"]; ok {
		t.Error("DRY_KEY should not have been written in dry-run mode")
	}
}
