package envfile

import (
	"fmt"
	"strings"
)

// SanitizeSummary describes the overall result of a Sanitize call.
type SanitizeSummary struct {
	Total   int
	Changed int
	Results []SanitizeResult
}

// SummarizeSanitize builds a SanitizeSummary from a result slice.
func SummarizeSanitize(results []SanitizeResult) SanitizeSummary {
	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
		}
	}
	return SanitizeSummary{
		Total:   len(results),
		Changed: changed,
		Results: results,
	}
}

// FormatSanitizeSummary returns a human-readable summary string.
func FormatSanitizeSummary(s SanitizeSummary) string {
	if s.Changed == 0 {
		return fmt.Sprintf("sanitize: %d entries checked, no changes needed", s.Total)
	}
	return fmt.Sprintf("sanitize: %d/%d entries modified", s.Changed, s.Total)
}

// FormatSanitizedEntries returns a detailed listing of changed entries.
// Sensitive values are masked when a Masker is provided (may be nil).
func FormatSanitizedEntries(results []SanitizeResult, m *Masker) string {
	var sb strings.Builder
	for _, r := range results {
		if !r.Changed {
			continue
		}
		before := r.Before
		after := r.After
		if m != nil && m.IsSensitive(r.Key) {
			before = m.MaskValue(r.Key, before)
			after = m.MaskValue(r.Key, after)
		}
		sb.WriteString(fmt.Sprintf("  ~ %s: %q -> %q\n", r.Key, before, after))
	}
	if sb.Len() == 0 {
		return "  (no changes)\n"
	}
	return sb.String()
}
