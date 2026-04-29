package envfile

import (
	"fmt"
	"sort"
)

// PinOptions controls how entries are pinned to specific values.
type PinOptions struct {
	// Overwrite allows pinning to replace an existing pinned value.
	Overwrite bool
	// FailOnMissing returns an error if a key to pin does not exist.
	FailOnMissing bool
}

// PinResult holds the outcome of a Pin operation.
type PinResult struct {
	Pinned  []string
	Skipped []string
	Missing []string
}

// Pin locks specific keys to given values in entries.
// pins is a map of key -> desired pinned value.
func Pin(entries []Entry, pins map[string]string, opts PinOptions) ([]Entry, PinResult, error) {
	result := PinResult{}
	index := make(map[string]int, len(entries))
	for i, e := range entries {
		index[e.Key] = i
	}

	out := make([]Entry, len(entries))
	copy(out, entries)

	// Process pins in deterministic order.
	keys := make([]string, 0, len(pins))
	for k := range pins {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := pins[k]
		i, exists := index[k]
		if !exists {
			if opts.FailOnMissing {
				return nil, result, fmt.Errorf("pin: key %q not found in entries", k)
			}
			result.Missing = append(result.Missing, k)
			continue
		}
		if out[i].Value == v && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		out[i].Value = v
		result.Pinned = append(result.Pinned, k)
	}

	return out, result, nil
}

// PinnedKeys returns the list of keys whose values match the provided pins map exactly.
func PinnedKeys(entries []Entry, pins map[string]string) []string {
	var matched []string
	for _, e := range entries {
		if v, ok := pins[e.Key]; ok && e.Value == v {
			matched = append(matched, e.Key)
		}
	}
	return matched
}
