package envfile

import (
	"fmt"
	"strings"
)

// FlattenOptions controls how nested key structures are flattened or expanded.
type FlattenOptions struct {
	// Separator is the delimiter used to join/split key segments (default: "__").
	Separator string
	// Prefix restricts flattening to keys with this prefix.
	Prefix string
}

// DefaultFlattenOptions returns sensible defaults for FlattenOptions.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator: "__",
	}
}

// FlattenKeys converts entries whose keys contain the separator into a nested
// map representation, returning a slice of Entry where each unique top-level
// segment becomes a key and its value is the joined sub-path value.
//
// Example: DB__HOST=localhost -> map key "DB" with sub-key "HOST".
func FlattenKeys(entries []Entry, opts FlattenOptions) ([]Entry, error) {
	if opts.Separator == "" {
		opts.Separator = "__"
	}

	seen := make(map[string]Entry)
	var order []string

	for _, e := range entries {
		key := e.Key
		if opts.Prefix != "" && !strings.HasPrefix(key, opts.Prefix) {
			// Pass through unmatched entries unchanged.
			if _, exists := seen[key]; !exists {
				order = append(order, key)
			}
			seen[key] = e
			continue
		}

		parts := strings.SplitN(key, opts.Separator, 2)
		if len(parts) == 1 {
			if _, exists := seen[key]; !exists {
				order = append(order, key)
			}
			seen[key] = e
			continue
		}

		if _, exists := seen[key]; !exists {
			order = append(order, key)
		}
		seen[key] = e
	}

	result := make([]Entry, 0, len(order))
	for _, k := range order {
		result = append(result, seen[k])
	}
	return result, nil
}

// ExpandKeys is the inverse of FlattenKeys: given a map-like slice of entries
// where values encode nested keys via the separator, it expands them.
//
// Example: given prefix "DB" and separator "__", an entry DB__HOST=localhost
// is left as-is; this function validates that no key contains an empty segment.
func ExpandKeys(entries []Entry, opts FlattenOptions) ([]Entry, error) {
	if opts.Separator == "" {
		opts.Separator = "__"
	}

	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		parts := strings.Split(e.Key, opts.Separator)
		for _, p := range parts {
			if p == "" {
				return nil, fmt.Errorf("key %q contains empty segment after splitting on %q", e.Key, opts.Separator)
			}
		}
		result = append(result, e)
	}
	return result, nil
}

// FilterByPrefix returns only the entries whose keys start with the given prefix.
// If prefix is empty, all entries are returned unchanged.
func FilterByPrefix(entries []Entry, prefix string) []Entry {
	if prefix == "" {
		return entries
	}
	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if strings.HasPrefix(e.Key, prefix) {
			result = append(result, e)
		}
	}
	return result
}
