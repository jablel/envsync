package envfile

import (
	"fmt"
	"strings"
)

// DiffSummary holds a human-readable summary of a diff result.
type DiffSummary struct {
	Added    int
	Removed  int
	Modified int
	Total    int
}

// SummarizeDiff returns a DiffSummary from a slice of DiffEntry.
func SummarizeDiff(entries []DiffEntry) DiffSummary {
	var s DiffSummary
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			s.Added++
		case StatusRemoved:
			s.Removed++
		case StatusModified:
			s.Modified++
		}
	}
	s.Total = s.Added + s.Removed + s.Modified
	return s
}

// FormatDiffSummary returns a short one-line summary string.
func FormatDiffSummary(s DiffSummary) string {
	if s.Total == 0 {
		return "no differences found"
	}
	parts := []string{}
	if s.Added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", s.Added))
	}
	if s.Removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", s.Removed))
	}
	if s.Modified > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", s.Modified))
	}
	return strings.Join(parts, ", ")
}

// FormatDiffEntries returns a multi-line human-readable diff listing.
// Sensitive values are masked when masker is non-nil.
func FormatDiffEntries(entries []DiffEntry, masker *Masker) string {
	if len(entries) == 0 {
		return "(no changes)"
	}
	var sb strings.Builder
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			val := maskIfNeeded(e.NewValue, e.Key, masker)
			fmt.Fprintf(&sb, "+ %s=%s\n", e.Key, val)
		case StatusRemoved:
			val := maskIfNeeded(e.OldValue, e.Key, masker)
			fmt.Fprintf(&sb, "- %s=%s\n", e.Key, val)
		case StatusModified:
			oldVal := maskIfNeeded(e.OldValue, e.Key, masker)
			newVal := maskIfNeeded(e.NewValue, e.Key, masker)
			fmt.Fprintf(&sb, "~ %s: %s -> %s\n", e.Key, oldVal, newVal)
		case StatusUnchanged:
			fmt.Fprintf(&sb, "  %s\n", e.Key)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func maskIfNeeded(val, key string, masker *Masker) string {
	if masker != nil && masker.IsSensitive(key) {
		return masker.MaskValue(val)
	}
	return val
}
