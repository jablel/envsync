package envfile

import (
	"fmt"
	"strings"
)

// AnnotateSummary holds statistics about an Annotate operation.
type AnnotateSummary struct {
	Applied []string
	Skipped []string
}

// SummarizeAnnotate builds an AnnotateSummary from a slice of AnnotateResult.
func SummarizeAnnotate(results []AnnotateResult) AnnotateSummary {
	s := AnnotateSummary{}
	for _, r := range results {
		switch {
		case r.Applied:
			s.Applied = append(s.Applied, r.Key)
		case r.Skipped:
			s.Skipped = append(s.Skipped, r.Key)
		}
	}
	return s
}

// FormatAnnotateSummary returns a human-readable summary string.
func FormatAnnotateSummary(s AnnotateSummary) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Annotate summary: %d applied, %d skipped\n",
		len(s.Applied), len(s.Skipped)))
	if len(s.Applied) > 0 {
		sb.WriteString("  Applied:\n")
		for _, k := range s.Applied {
			sb.WriteString(fmt.Sprintf("    + %s\n", k))
		}
	}
	if len(s.Skipped) > 0 {
		sb.WriteString("  Skipped (already annotated):\n")
		for _, k := range s.Skipped {
			sb.WriteString(fmt.Sprintf("    ~ %s\n", k))
		}
	}
	return sb.String()
}

// FormatAnnotatedEntries renders entries with their annotations for display.
func FormatAnnotatedEntries(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		if e.Comment != "" {
			sb.WriteString(fmt.Sprintf("%-20s = %-20s  # %s\n", e.Key, e.Value, e.Comment))
		} else {
			sb.WriteString(fmt.Sprintf("%-20s = %s\n", e.Key, e.Value))
		}
	}
	return sb.String()
}
