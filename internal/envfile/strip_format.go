package envfile

import (
	"fmt"
	"strings"
)

// StripSummary returns a human-readable summary of a StripResult.
func StripSummary(result StripResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Strip summary: %d kept, %d removed\n",
		len(result.Kept), len(result.Removed)))
	if len(result.Removed) > 0 {
		sb.WriteString("Removed keys:\n")
		for _, e := range result.Removed {
			sb.WriteString(fmt.Sprintf("  - %s\n", e.Key))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatStrippedEntries formats the kept entries after stripping, optionally
// masking sensitive values.
func FormatStrippedEntries(result StripResult, mask bool) string {
	if len(result.Kept) == 0 {
		return "(no entries remaining)"
	}
	masker := NewMasker()
	var sb strings.Builder
	for _, e := range result.Kept {
		val := e.Value
		if mask && masker.IsSensitive(e.Key) {
			val = masker.MaskValue(val)
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, val))
	}
	return strings.TrimRight(sb.String(), "\n")
}
