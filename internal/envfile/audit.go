package envfile

import (
	"fmt"
	"time"
)

// AuditAction represents the type of change recorded in an audit log.
type AuditAction string

const (
	AuditAdded    AuditAction = "added"
	AuditRemoved  AuditAction = "removed"
	AuditModified AuditAction = "modified"
	AuditMasked   AuditAction = "masked"
)

// AuditEntry records a single change event for an environment variable.
type AuditEntry struct {
	Timestamp time.Time
	Action    AuditAction
	Key       string
	OldValue  string
	NewValue  string
}

// AuditLog holds a collection of audit entries.
type AuditLog struct {
	Entries []AuditEntry
	masker  *Masker
}

// NewAuditLog creates a new AuditLog, optionally using a Masker to redact values.
func NewAuditLog(m *Masker) *AuditLog {
	return &AuditLog{masker: m}
}

// Record appends a new audit entry with the current timestamp.
func (a *AuditLog) Record(action AuditAction, key, oldVal, newVal string) {
	if a.masker != nil && a.masker.IsSensitive(key) {
		oldVal = a.masker.MaskValue(oldVal)
		newVal = a.masker.MaskValue(newVal)
		action = AuditMasked
	}
	a.Entries = append(a.Entries, AuditEntry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Key:       key,
		OldValue:  oldVal,
		NewValue:  newVal,
	})
}

// FromDiff populates an AuditLog from a slice of DiffEntry values.
func (a *AuditLog) FromDiff(diffs []DiffEntry) {
	for _, d := range diffs {
		switch d.Status {
		case StatusAdded:
			a.Record(AuditAdded, d.Key, "", d.NewValue)
		case StatusRemoved:
			a.Record(AuditRemoved, d.Key, d.OldValue, "")
		case StatusModified:
			a.Record(AuditModified, d.Key, d.OldValue, d.NewValue)
		}
	}
}

// Summary returns a human-readable summary of the audit log.
func (a *AuditLog) Summary() string {
	if len(a.Entries) == 0 {
		return "no changes recorded"
	}
	s := fmt.Sprintf("%d change(s) recorded:\n", len(a.Entries))
	for _, e := range a.Entries {
		s += fmt.Sprintf("  [%s] %s: %q -> %q\n",
			e.Action, e.Key, e.OldValue, e.NewValue)
	}
	return s
}
