package envfile

import (
	"fmt"
	"strings"
)

// TransformSummary describes the result of a Transform call.
type TransformSummary struct {
	Total   int
	Changed int
	Skipped int // entries left unchanged due to SkipErrors
}

// SummarizeTransform compares before and after slices and returns a
// TransformSummary. Both slices must have the same length.
func SummarizeTransform(before, after []Entry) TransformSummary {
	s := TransformSummary{Total: len(before)}
	for i := range before {
		if i >= len(after) {
			break
		}
		if before[i].Key != after[i].Key || before[i].Value != after[i].Value {
			s.Changed++
		}
	}
	return s
}

// FormatTransformSummary returns a human-readable summary string.
func FormatTransformSummary(s TransformSummary) string {
	var b strings.Builder
	fmt.Fprintf(&b, "transform: %d entries processed", s.Total)
	if s.Changed > 0 {
		fmt.Fprintf(&b, ", %d changed", s.Changed)
	}
	if s.Skipped > 0 {
		fmt.Fprintf(&b, ", %d skipped (errors)", s.Skipped)
	}
	return b.String()
}

// FormatTransformedEntries renders the transformed entries as a diff-style
// summary showing only the keys whose values changed.
func FormatTransformedEntries(before, after []Entry, masker *Masker) string {
	var b strings.Builder
	for i := range before {
		if i >= len(after) {
			break
		}
		prev := before[i]
		next := after[i]
		if prev.Key == next.Key && prev.Value == next.Value {
			continue
		}
		prevVal := prev.Value
		nextVal := next.Value
		if masker != nil && masker.IsSensitive(next.Key) {
			prevVal = masker.MaskValue(prevVal)
			nextVal = masker.MaskValue(nextVal)
		}
		fmt.Fprintf(&b, "~ %s: %q -> %q\n", next.Key, prevVal, nextVal)
	}
	if b.Len() == 0 {
		return "(no changes)\n"
	}
	return b.String()
}
