package envfile

import (
	"fmt"
	"strings"
)

// PinSummary returns a human-readable summary of a PinResult.
func PinSummary(res PinResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Pinned: %d, Skipped: %d, Missing: %d\n",
		len(res.Pinned), len(res.Skipped), len(res.Missing)))
	if len(res.Pinned) > 0 {
		sb.WriteString("  Pinned keys:\n")
		for _, k := range res.Pinned {
			sb.WriteString(fmt.Sprintf("    + %s\n", k))
		}
	}
	if len(res.Skipped) > 0 {
		sb.WriteString("  Skipped (already pinned):\n")
		for _, k := range res.Skipped {
			sb.WriteString(fmt.Sprintf("    ~ %s\n", k))
		}
	}
	if len(res.Missing) > 0 {
		sb.WriteString("  Missing keys:\n")
		for _, k := range res.Missing {
			sb.WriteString(fmt.Sprintf("    ? %s\n", k))
		}
	}
	return sb.String()
}

// FormatPinnedEntries formats entries highlighting which keys are pinned.
func FormatPinnedEntries(entries []Entry, pins map[string]string, masker *Masker) string {
	var sb strings.Builder
	for _, e := range entries {
		val := e.Value
		if masker != nil && masker.IsSensitive(e.Key) {
			val = masker.MaskValue(val)
		}
		if _, pinned := pins[e.Key]; pinned {
			sb.WriteString(fmt.Sprintf("[PINNED] %s=%s\n", e.Key, val))
		} else {
			sb.WriteString(fmt.Sprintf("         %s=%s\n", e.Key, val))
		}
	}
	return sb.String()
}
