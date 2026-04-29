package envfile

import (
	"fmt"
	"strings"
)

// ResolveSummary returns a human-readable summary of a ResolveResult.
func ResolveSummary(r ResolveResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Resolved: %d applied", r.Applied)
	if len(r.Conflicts) > 0 {
		fmt.Fprintf(&sb, ", %d conflict(s) kept base", len(r.Conflicts))
	}
	fmt.Fprintf(&sb, ", %d total entries\n", len(r.Entries))
	return sb.String()
}

// FormatResolvedEntries formats the resolved entries with change markers.
// base is used to detect which values were overridden.
func FormatResolvedEntries(r ResolveResult, base []Entry) string {
	baseMap := make(map[string]string, len(base))
	for _, e := range base {
		baseMap[e.Key] = e.Value
	}

	var sb strings.Builder
	for _, e := range r.Entries {
		origVal, existed := baseMap[e.Key]
		switch {
		case !existed:
			fmt.Fprintf(&sb, "+ %s=%s\n", e.Key, e.Value)
		case existed && origVal != e.Value:
			fmt.Fprintf(&sb, "~ %s=%s (was: %s)\n", e.Key, e.Value, origVal)
		default:
			fmt.Fprintf(&sb, "  %s=%s\n", e.Key, e.Value)
		}
	}
	return sb.String()
}

// FormatConflicts returns a formatted list of conflicting keys.
func FormatConflicts(conflicts []string) string {
	if len(conflicts) == 0 {
		return "No conflicts.\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Conflicts (%d keys kept base value):\n", len(conflicts))
	for _, k := range conflicts {
		fmt.Fprintf(&sb, "  - %s\n", k)
	}
	return sb.String()
}
