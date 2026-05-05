package envfile

import (
	"fmt"
	"strings"
)

// ChainSummary returns a one-line summary of a ChainResult.
func ChainSummary(r ChainResult) string {
	total := len(r.Steps)
	applied := 0
	failed := 0
	for _, s := range r.Steps {
		if s.Applied {
			applied++
		}
		if s.Err != nil {
			failed++
		}
	}
	if failed == 0 {
		return fmt.Sprintf("chain: %d/%d steps applied, %d entries", applied, total, len(r.Entries))
	}
	return fmt.Sprintf("chain: %d/%d steps applied (%d failed), %d entries", applied, total, failed, len(r.Entries))
}

// FormatChainSteps returns a multi-line report of each step's outcome.
func FormatChainSteps(r ChainResult) string {
	if len(r.Steps) == 0 {
		return "  (no steps)"
	}
	var sb strings.Builder
	for _, s := range r.Steps {
		if s.Err != nil {
			fmt.Fprintf(&sb, "  [FAIL] %-30s  error: %v\n", s.Name, s.Err)
		} else {
			fmt.Fprintf(&sb, "  [OK]   %s\n", s.Name)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
