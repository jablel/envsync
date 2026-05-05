package envfile

import (
	"fmt"
	"strings"
)

// ReorderSummary returns a human-readable summary of a ReorderResult.
func ReorderSummary(r ReorderResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "reorder: %d entries", len(r.Entries))
	if len(r.Moved) > 0 {
		fmt.Fprintf(&sb, ", %d moved", len(r.Moved))
	}
	if len(r.Missing) > 0 {
		fmt.Fprintf(&sb, ", %d missing", len(r.Missing))
	}
	return sb.String()
}

// FormatReorderDiff renders a before/after view of moved keys.
func FormatReorderDiff(original, result []Entry) string {
	if len(result) == 0 {
		return "(no entries)"
	}

	origPos := make(map[string]int, len(original))
	for i, e := range original {
		origPos[e.Key] = i
	}

	var sb strings.Builder
	for newIdx, e := range result {
		old, known := origPos[e.Key]
		if known && old != newIdx {
			fmt.Fprintf(&sb, "  ~ %s  (was %d, now %d)\n", e.Key, old+1, newIdx+1)
		} else {
			fmt.Fprintf(&sb, "    %s\n", e.Key)
		}
	}
	return sb.String()
}

// FormatMissingReorderKeys formats keys that were listed in the order but
// absent from the source entries.
func FormatMissingReorderKeys(missing []string) string {
	if len(missing) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("missing keys (listed in order but not found):\n")
	for _, k := range missing {
		fmt.Fprintf(&sb, "  - %s\n", k)
	}
	return sb.String()
}
