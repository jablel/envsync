package envfile

import (
	"fmt"
	"strings"
)

// FreezeSummary returns a human-readable summary of a FreezeResult.
func FreezeSummary(res FreezeResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("frozen: %d", len(res.Frozen)))
	if len(res.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf(", skipped (already frozen): %d", len(res.Skipped)))
	}
	if len(res.Errors) > 0 {
		sb.WriteString(fmt.Sprintf(", errors: %d", len(res.Errors)))
	}
	return sb.String()
}

// FormatFrozenEntries returns a formatted list of frozen keys for display.
func FormatFrozenEntries(entries []Entry, masker *Masker) string {
	var sb strings.Builder
	for _, e := range entries {
		if !IsFrozen(e) {
			continue
		}
		val := e.Value
		if masker != nil && masker.IsSensitive(e.Key) {
			val = masker.MaskValue(e.Key, val)
		}
		sb.WriteString(fmt.Sprintf("  [frozen] %s=%s\n", e.Key, val))
	}
	if sb.Len() == 0 {
		return "  (no frozen entries)\n"
	}
	return sb.String()
}
