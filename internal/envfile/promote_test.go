package envfile

import (
	"testing"
)

func srcEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "prod-db.example.com"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
	}
}

func dstEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "APP_ENV", Value: "staging"},
	}
}

func TestPromote_NewKeysAdded(t *testing.T) {
	src := srcEntries()
	dst := dstEntries()
	updated, res, err := Promote(src, dst, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	m := entriesToMap(updated)
	if m["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", m["DB_PORT"])
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := srcEntries()
	dst := dstEntries()
	updated, res, err := Promote(src, dst, PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten, got %d", len(res.Overwritten))
	}
	m := entriesToMap(updated)
	if m["DB_HOST"] != "prod-db.example.com" {
		t.Errorf("expected overwritten DB_HOST, got %q", m["DB_HOST"])
	}
}

func TestPromote_FilterKeys(t *testing.T) {
	src := srcEntries()
	dst := dstEntries()
	_, res, err := Promote(src, dst, PromoteOptions{Keys: []string{"DB_PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 1 || res.Promoted[0] != "DB_PORT" {
		t.Errorf("expected only DB_PORT promoted, got %v", res.Promoted)
	}
}

func TestPromote_MissingFilterKey(t *testing.T) {
	src := srcEntries()
	dst := dstEntries()
	_, _, err := Promote(src, dst, PromoteOptions{Keys: []string{"NONEXISTENT"}})
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestPromote_DryRunDoesNotModify(t *testing.T) {
	src := srcEntries()
	dst := dstEntries()
	_, res, err := Promote(src, dst, PromoteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dst) != 2 {
		t.Errorf("dry run must not modify dst, got len=%d", len(dst))
	}
	if len(res.Promoted) == 0 {
		t.Error("expected promoted keys reported even in dry run")
	}
}

func entriesToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
