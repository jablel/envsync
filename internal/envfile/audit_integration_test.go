package envfile_test

import (
	"strings"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func TestAudit_DiffAndSummaryRoundTrip(t *testing.T) {
	base := []envfile.Entry{
		{Key: "APP_ENV", Value: "staging"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "OLD_VAR", Value: "remove-me"},
	}
	override := []envfile.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "NEW_VAR", Value: "hello"},
	}

	diffs := envfile.Diff(base, override)

	log := envfile.NewAuditLog(nil)
	log.FromDiff(diffs)

	if len(log.Entries) != 3 {
		t.Fatalf("expected 3 audit entries, got %d", len(log.Entries))
	}

	summary := log.Summary()
	for _, want := range []string{"APP_ENV", "OLD_VAR", "NEW_VAR"} {
		if !strings.Contains(summary, want) {
			t.Errorf("summary missing key %q", want)
		}
	}
}

func TestAudit_MasksSecretsInDiff(t *testing.T) {
	base := []envfile.Entry{
		{Key: "API_SECRET", Value: "old-token-abc"},
	}
	override := []envfile.Entry{
		{Key: "API_SECRET", Value: "new-token-xyz"},
	}

	diffs := envfile.Diff(base, override)

	m := envfile.NewMasker(nil)
	log := envfile.NewAuditLog(m)
	log.FromDiff(diffs)

	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if strings.Contains(e.OldValue, "old-token-abc") {
		t.Error("old secret value should be masked in audit log")
	}
	if strings.Contains(e.NewValue, "new-token-xyz") {
		t.Error("new secret value should be masked in audit log")
	}
}
