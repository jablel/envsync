package envfile

import (
	"fmt"
	"strings"
)

// DiffSnapshotSummary holds counts of changes between a snapshot and current state.
type DiffSnapshotSummary struct {
	Added    int
	Removed  int
	Modified int
	Unchanged int
}

// SummarizeDiffSnapshot produces a summary from a DiffResult.
func SummarizeDiffSnapshot(d DiffResult) DiffSnapshotSummary {
	return DiffSnapshotSummary{
		Added:     len(d.Added),
		Removed:   len(d.Removed),
		Modified:  len(d.Modified),
		Unchanged: len(d.Unchanged),
	}
}

// FormatDiffSnapshotSummary returns a human-readable summary string.
func FormatDiffSnapshotSummary(s DiffSnapshotSummary) string {
	if s.Added == 0 && s.Removed == 0 && s.Modified == 0 {
		return fmt.Sprintf("No changes vs snapshot (%d unchanged)", s.Unchanged)
	}
	parts := []string{}
	if s.Added > 0 {
		parts = append(parts, fmt.Sprintf("+%d added", s.Added))
	}
	if s.Removed > 0 {
		parts = append(parts, fmt.Sprintf("-%d removed", s.Removed))
	}
	if s.Modified > 0 {
		parts = append(parts, fmt.Sprintf("~%d modified", s.Modified))
	}
	return fmt.Sprintf("Snapshot diff: %s (%d unchanged)", strings.Join(parts, ", "), s.Unchanged)
}

// FormatDiffSnapshotEntries renders added/removed/modified entries with
// snapshot-aware prefixes. Sensitive values are masked when masker != nil.
func FormatDiffSnapshotEntries(d DiffResult, masker *Masker) string {
	var sb strings.Builder
	for _, e := range d.Added {
		v := e.Value
		if masker != nil && masker.IsSensitive(e.Key) {
			v = masker.MaskValue(v)
		}
		fmt.Fprintf(&sb, "[snapshot] + %s=%s\n", e.Key, v)
	}
	for _, c := range d.Modified {
		oldV, newV := c.Old.Value, c.New.Value
		if masker != nil && masker.IsSensitive(c.Old.Key) {
			oldV = masker.MaskValue(oldV)
			newV = masker.MaskValue(newV)
		}
		fmt.Fprintf(&sb, "[snapshot] ~ %s: %s -> %s\n", c.Old.Key, oldV, newV)
	}
	for _, e := range d.Removed {
		v := e.Value
		if masker != nil && masker.IsSensitive(e.Key) {
			v = masker.MaskValue(v)
		}
		fmt.Fprintf(&sb, "[snapshot] - %s=%s\n", e.Key, v)
	}
	return sb.String()
}
