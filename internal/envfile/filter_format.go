package envfile

import (
	"fmt"
	"strings"
)

// FilterSummary holds statistics about a filter operation.
type FilterSummary struct {
	Total    int
	Matched  int
	Excluded int
}

// SummarizeFilter computes a FilterSummary given the original and filtered slices.
func SummarizeFilter(original, filtered []Entry) FilterSummary {
	return FilterSummary{
		Total:    len(original),
		Matched:  len(filtered),
		Excluded: len(original) - len(filtered),
	}
}

// FormatFilterSummary returns a human-readable summary string.
func FormatFilterSummary(s FilterSummary) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Filter results: %d matched / %d total", s.Matched, s.Total))
	if s.Excluded > 0 {
		sb.WriteString(fmt.Sprintf(" (%d excluded)", s.Excluded))
	}
	return sb.String()
}

// FormatFilteredEntries renders filtered entries as key=value lines,
// masking sensitive values when a Masker is provided.
func FormatFilteredEntries(entries []Entry, masker *Masker) string {
	if len(entries) == 0 {
		return "(no entries matched)"
	}
	var sb strings.Builder
	for _, e := range entries {
		val := e.Value
		if masker != nil && masker.IsSensitive(e.Key) {
			val = masker.MaskValue(e.Key, val)
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, val))
	}
	return sb.String()
}
