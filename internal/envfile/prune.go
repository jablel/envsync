package envfile

import "strings"

// PruneOptions controls which entries are removed during pruning.
type PruneOptions struct {
	// RemoveCommented removes entries whose values start with a comment marker.
	RemoveCommented bool
	// RemovePlaceholders removes entries whose values match common placeholder patterns.
	RemovePlaceholders bool
	// RemoveKeys removes entries with these exact keys.
	RemoveKeys []string
	// RemovePrefixes removes entries whose keys start with any of these prefixes.
	RemovePrefixes []string
}

// PruneResult holds the outcome of a Prune operation.
type PruneResult struct {
	Kept    []Entry
	Removed []Entry
}

var defaultPlaceholders = []string{
	"TODO", "FIXME", "CHANGEME", "<your-", "your_", "example", "placeholder",
}

// Prune removes entries from the list based on the given options.
func Prune(entries []Entry, opts PruneOptions) PruneResult {
	removeKeySet := make(map[string]struct{}, len(opts.RemoveKeys))
	for _, k := range opts.RemoveKeys {
		removeKeySet[k] = struct{}{}
	}

	var result PruneResult
	for _, e := range entries {
		if shouldPrune(e, opts, removeKeySet) {
			result.Removed = append(result.Removed, e)
		} else {
			result.Kept = append(result.Kept, e)
		}
	}
	return result
}

// PrunedKeys returns only the keys that were removed.
func PrunedKeys(result PruneResult) []string {
	keys := make([]string, 0, len(result.Removed))
	for _, e := range result.Removed {
		keys = append(keys, e.Key)
	}
	return keys
}

func shouldPrune(e Entry, opts PruneOptions, removeKeySet map[string]struct{}) bool {
	if _, ok := removeKeySet[e.Key]; ok {
		return true
	}
	for _, prefix := range opts.RemovePrefixes {
		if strings.HasPrefix(e.Key, prefix) {
			return true
		}
	}
	if opts.RemovePlaceholders {
		lower := strings.ToLower(e.Value)
		for _, ph := range defaultPlaceholders {
			if strings.Contains(lower, strings.ToLower(ph)) {
				return true
			}
		}
	}
	if opts.RemoveCommented && strings.HasPrefix(strings.TrimSpace(e.Value), "#") {
		return true
	}
	return false
}
