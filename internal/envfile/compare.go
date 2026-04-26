package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// CompareResult holds the result of comparing two env file sets.
type CompareResult struct {
	Matching    []string
	Mismatched  map[string][2]string // key -> [base, other]
	OnlyInBase  []string
	OnlyInOther []string
}

// Summary returns a human-readable summary of the comparison.
func (r *CompareResult) Summary() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Matching keys:     %d\n", len(r.Matching))
	fmt.Fprintf(&sb, "Mismatched values: %d\n", len(r.Mismatched))
	fmt.Fprintf(&sb, "Only in base:      %d\n", len(r.OnlyInBase))
	fmt.Fprintf(&sb, "Only in other:     %d\n", len(r.OnlyInOther))
	return sb.String()
}

// IsIdentical returns true when both env files are fully equivalent.
func (r *CompareResult) IsIdentical() bool {
	return len(r.Mismatched) == 0 &&
		len(r.OnlyInBase) == 0 &&
		len(r.OnlyInOther) == 0
}

// Compare performs a detailed comparison between a base and another set of
// env entries, returning a CompareResult that describes every category of
// difference.
func Compare(base, other []Entry) *CompareResult {
	baseMap := make(map[string]string, len(base))
	for _, e := range base {
		baseMap[e.Key] = e.Value
	}

	otherMap := make(map[string]string, len(other))
	for _, e := range other {
		otherMap[e.Key] = e.Value
	}

	result := &CompareResult{
		Mismatched: make(map[string][2]string),
	}

	for k, bv := range baseMap {
		ov, exists := otherMap[k]
		if !exists {
			result.OnlyInBase = append(result.OnlyInBase, k)
		} else if bv == ov {
			result.Matching = append(result.Matching, k)
		} else {
			result.Mismatched[k] = [2]string{bv, ov}
		}
	}

	for k := range otherMap {
		if _, exists := baseMap[k]; !exists {
			result.OnlyInOther = append(result.OnlyInOther, k)
		}
	}

	sort.Strings(result.Matching)
	sort.Strings(result.OnlyInBase)
	sort.Strings(result.OnlyInOther)

	return result
}
