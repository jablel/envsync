package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatOptions controls how env entries are rendered to text.
type FormatOptions struct {
	SortKeys    bool
	MaskSecrets bool
	ShowDiff    bool
}

// Format renders a slice of Entry values into .env file text.
func Format(entries []Entry, opts FormatOptions) string {
	if opts.SortKeys {
		sorted := make([]Entry, len(entries))
		copy(sorted, entries)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Key < sorted[j].Key
		})
		entries = sorted
	}

	var sb strings.Builder
	for _, e := range entries {
		val := e.Value
		if opts.MaskSecrets {
			masker := NewMasker(nil)
			if masker.IsSensitive(e.Key) {
				val = masker.MaskValue(val)
			}
		}
		if needsQuoting(val) {
			sb.WriteString(fmt.Sprintf("%s=%q\n", e.Key, val))
		} else {
			sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, val))
		}
	}
	return sb.String()
}

// FormatDiff renders a DiffResult as a human-readable diff string.
func FormatDiff(result DiffResult, mask bool) string {
	var sb strings.Builder
	masker := NewMasker(nil)

	for _, e := range result.Added {
		val := e.Value
		if mask && masker.IsSensitive(e.Key) {
			val = masker.MaskValue(val)
		}
		sb.WriteString(fmt.Sprintf("+ %s=%s\n", e.Key, val))
	}

	for _, e := range result.Removed {
		val := e.Value
		if mask && masker.IsSensitive(e.Key) {
			val = masker.MaskValue(val)
		}
		sb.WriteString(fmt.Sprintf("- %s=%s\n", e.Key, val))
	}

	for _, c := range result.Modified {
		oldVal, newVal := c.OldValue, c.NewValue
		if mask && masker.IsSensitive(c.Key) {
			oldVal = masker.MaskValue(oldVal)
			newVal = masker.MaskValue(newVal)
		}
		sb.WriteString(fmt.Sprintf("~ %s: %s -> %s\n", c.Key, oldVal, newVal))
	}

	return sb.String()
}

// needsQuoting returns true if the value contains characters that require quoting.
func needsQuoting(val string) bool {
	return strings.ContainsAny(val, " \t#\n")
}
