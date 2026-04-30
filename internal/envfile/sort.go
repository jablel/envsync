package envfile

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering strategy for env entries.
type SortOrder int

const (
	SortAlpha      SortOrder = iota // alphabetical by key
	SortAlphaDesc                   // reverse alphabetical
	SortByPrefix                    // group by prefix (e.g. DB_, AWS_), then alpha within group
	SortSecretLast                  // non-sensitive keys first, sensitive keys last
)

// SortOptions configures how entries are sorted.
type SortOptions struct {
	Order      SortOrder
	Masker     *Masker // required for SortSecretLast
	StableOrig bool    // if true, preserve relative order for equal keys
}

// DefaultSortOptions returns sensible defaults.
func DefaultSortOptions() SortOptions {
	return SortOptions{
		Order: SortAlpha,
	}
}

// Sort returns a new slice of entries sorted according to opts.
func Sort(entries []Entry, opts SortOptions) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)

	switch opts.Order {
	case SortAlphaDesc:
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].Key > out[j].Key
		})
	case SortByPrefix:
		sort.SliceStable(out, func(i, j int) bool {
			pi := extractPrefix(out[i].Key)
			pj := extractPrefix(out[j].Key)
			if pi != pj {
				return pi < pj
			}
			return out[i].Key < out[j].Key
		})
	case SortSecretLast:
		m := opts.Masker
		if m == nil {
			m = NewMasker(nil)
		}
		sort.SliceStable(out, func(i, j int) bool {
			si := m.IsSensitive(out[i].Key)
			sj := m.IsSensitive(out[j].Key)
			if si != sj {
				return !si // non-sensitive first
			}
			return out[i].Key < out[j].Key
		})
	default: // SortAlpha
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].Key < out[j].Key
		})
	}

	return out
}

// SortedKeys returns only the keys from entries in sorted order.
func SortedKeys(entries []Entry, opts SortOptions) []string {
	sorted := Sort(entries, opts)
	keys := make([]string, len(sorted))
	for i, e := range sorted {
		keys[i] = e.Key
	}
	return keys
}

// extractPrefix returns the prefix portion of a key (e.g. "DB" from "DB_HOST").
// Keys without an underscore return the full key as the prefix.
func extractPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
