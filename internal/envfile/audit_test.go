package envfile

import (
	"testing"
)

func TestAuditLog_RecordBasic(t *testing.T) {
	log := NewAuditLog(nil)
	log.Record(AuditAdded, "APP_ENV", "", "production")

	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if e.Action != AuditAdded {
		t.Errorf("expected action %q, got %q", AuditAdded, e.Action)
	}
	if e.Key != "APP_ENV" {
		t.Errorf("unexpected key: %s", e.Key)
	}
	if e.NewValue != "production" {
		t.Errorf("unexpected new value: %s", e.NewValue)
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestAuditLog_MasksSensitiveKeys(t *testing.T) {
	m := NewMasker(nil)
	log := NewAuditLog(m)
	log.Record(AuditModified, "DB_PASSWORD", "old-secret", "new-secret")

	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if e.Action != AuditMasked {
		t.Errorf("expected action %q, got %q", AuditMasked, e.Action)
	}
	if e.OldValue == "old-secret" {
		t.Error("old value should be masked")
	}
	if e.NewValue == "new-secret" {
		t.Error("new value should be masked")
	}
}

func TestAuditLog_FromDiff(t *testing.T) {
	diffs := []DiffEntry{
		{Key: "NEW_KEY", Status: StatusAdded, NewValue: "val1"},
		{Key: "OLD_KEY", Status: StatusRemoved, OldValue: "val2"},
		{Key: "MOD_KEY", Status: StatusModified, OldValue: "old", NewValue: "new"},
		{Key: "SAME_KEY", Status: StatusUnchanged, OldValue: "same", NewValue: "same"},
	}

	log := NewAuditLog(nil)
	log.FromDiff(diffs)

	if len(log.Entries) != 3 {
		t.Fatalf("expected 3 entries (unchanged skipped), got %d", len(log.Entries))
	}
	if log.Entries[0].Action != AuditAdded {
		t.Errorf("first entry should be added, got %q", log.Entries[0].Action)
	}
	if log.Entries[1].Action != AuditRemoved {
		t.Errorf("second entry should be removed, got %q", log.Entries[1].Action)
	}
	if log.Entries[2].Action != AuditModified {
		t.Errorf("third entry should be modified, got %q", log.Entries[2].Action)
	}
}

func TestAuditLog_Summary_Empty(t *testing.T) {
	log := NewAuditLog(nil)
	got := log.Summary()
	if got != "no changes recorded" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestAuditLog_Summary_WithEntries(t *testing.T) {
	log := NewAuditLog(nil)
	log.Record(AuditAdded, "FOO", "", "bar")
	log.Record(AuditRemoved, "BAZ", "qux", "")

	summary := log.Summary()
	if summary == "no changes recorded" {
		t.Error("expected non-empty summary")
	}
	for _, want := range []string{"FOO", "BAZ", "added", "removed"} {
		if !containsStr(summary, want) {
			t.Errorf("summary missing %q", want)
		}
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsStr(s[1:], substr) ||
		len(s) >= len(substr) && s[:len(substr)] == substr)
}
