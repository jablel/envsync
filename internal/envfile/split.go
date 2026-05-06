package envfile

import "fmt"

// SplitOptions controls how entries are split into multiple groups.
type SplitOptions struct {
	// Keys defines explicit key lists per output group name.
	Keys map[string][]string
	// Prefixes maps group names to key prefixes.
	Prefixes map[string][]string
	// DropUnmatched silently drops entries that don't belong to any group.
	DropUnmatched bool
}

// SplitResult holds the output groups and any unmatched entries.
type SplitResult struct {
	Groups    map[string][]Entry
	Unmatched []Entry
}

// Split partitions entries into named groups based on explicit keys or prefixes.
// Entries may appear in multiple groups if they match more than one rule.
func Split(entries []Entry, opts SplitOptions) (SplitResult, error) {
	if len(opts.Keys) == 0 && len(opts.Prefixes) == 0 {
		return SplitResult{}, fmt.Errorf("split: at least one key list or prefix group must be specified")
	}

	groups := make(map[string][]Entry)

	// Pre-build key sets for O(1) lookup.
	keySets := make(map[string]map[string]struct{}, len(opts.Keys))
	for group, keys := range opts.Keys {
		set := make(map[string]struct{}, len(keys))
		for _, k := range keys {
			set[k] = struct{}{}
		}
		keySets[group] = set
	}

	var unmatched []Entry

	for _, e := range entries {
		matched := false

		for group, set := range keySets {
			if _, ok := set[e.Key]; ok {
				groups[group] = append(groups[group], e)
				matched = true
			}
		}

		for group, prefixes := range opts.Prefixes {
			for _, p := range prefixes {
				if hasPrefix(e.Key, p) {
					groups[group] = append(groups[group], e)
					matched = true
					break
				}
			}
		}

		if !matched {
			unmatched = append(unmatched, e)
		}
	}

	return SplitResult{
		Groups:    groups,
		Unmatched: unmatched,
	}, nil
}

// SplitGroupNames returns sorted group names present in a SplitResult.
func SplitGroupNames(result SplitResult) []string {
	names := make([]string, 0, len(result.Groups))
	for name := range result.Groups {
		names = append(names, name)
	}
	sortStrings(names)
	return names
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
