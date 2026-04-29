package envfile

import (
	"regexp"
	"strings"
)

// FilterOptions controls how entries are filtered.
type FilterOptions struct {
	// Keys is an explicit list of keys to include (empty = all).
	Keys []string
	// Prefix filters entries whose keys start with the given prefix.
	Prefix string
	// Pattern is a regex applied to key names.
	Pattern string
	// ExcludeSecrets removes sensitive entries when true.
	ExcludeSecrets bool
	// Masker is used when ExcludeSecrets is true.
	Masker *Masker
}

// Filter returns a subset of entries matching the provided options.
// Multiple criteria are ANDed together.
func Filter(entries []Entry, opts FilterOptions) ([]Entry, error) {
	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
	}

	var result []Entry
	for _, e := range entries {
		if len(keySet) > 0 {
			if _, ok := keySet[e.Key]; !ok {
				continue
			}
		}
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if re != nil && !re.MatchString(e.Key) {
			continue
		}
		if opts.ExcludeSecrets && opts.Masker != nil && opts.Masker.IsSensitive(e.Key) {
			continue
		}
		result = append(result, e)
	}
	return result, nil
}

// FilterKeys returns only the key names from the filtered result.
func FilterKeys(entries []Entry, opts FilterOptions) ([]string, error) {
	filtered, err := Filter(entries, opts)
	if err != nil {
		return nil, err
	}
	keys := make([]string, len(filtered))
	for i, e := range filtered {
		keys[i] = e.Key
	}
	return keys, nil
}
