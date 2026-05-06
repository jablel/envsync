package envfile

import (
	"fmt"
	"strings"
)

// TagSummary holds aggregate statistics from a Tag operation.
type TagSummary struct {
	TotalTagged int
	TotalSkipped int
	Results []TagResult
}

// SummarizeTag builds a TagSummary from a slice of TagResults.
func SummarizeTag(results []TagResult) TagSummary {
	s := TagSummary{Results: results}
	for _, r := range results {
		s.TotalTagged += len(r.Tagged)
		s.TotalSkipped += len(r.Skipped)
	}
	return s
}

// FormatTagSummary returns a human-readable summary string.
func FormatTagSummary(s TagSummary) string {
	if s.TotalTagged == 0 && s.TotalSkipped == 0 {
		return "tag: no entries matched"
	}
	return fmt.Sprintf("tag: %d tag(s) applied, %d skipped", s.TotalTagged, s.TotalSkipped)
}

// FormatTaggedEntries returns a detailed per-entry report.
func FormatTaggedEntries(results []TagResult) string {
	if len(results) == 0 {
		return "(no changes)"
	}
	var sb strings.Builder
	for _, r := range results {
		if len(r.Tagged) > 0 {
			sb.WriteString(fmt.Sprintf("  + %-24s tagged: %s\n", r.Key, strings.Join(r.Tagged, ", ")))
		}
		if len(r.Skipped) > 0 {
			sb.WriteString(fmt.Sprintf("  ~ %-24s skipped: %s\n", r.Key, strings.Join(r.Skipped, ", ")))
		}
	}
	return sb.String()
}
