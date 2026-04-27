package envfile

import "sort"

// GroupOptions controls how entries are grouped.
type GroupOptions struct {
	// Prefixes defines known prefix groups (e.g. "DB", "AWS").
	// Entries whose keys start with PREFIX_ are placed in that group.
	Prefixes []string
	// DefaultGroup is the name for entries that don't match any prefix.
	DefaultGroup string
}

// DefaultGroupOptions returns sensible defaults.
func DefaultGroupOptions() GroupOptions {
	return GroupOptions{
		Prefixes:     []string{},
		DefaultGroup: "general",
	}
}

// GroupedEntries maps group name -> entries belonging to that group.
type GroupedEntries map[string][]Entry

// Group partitions entries by key prefix.
// Entries are matched against opts.Prefixes in order; the first match wins.
// Unmatched entries fall into opts.DefaultGroup.
func Group(entries []Entry, opts GroupOptions) GroupedEntries {
	if opts.DefaultGroup == "" {
		opts.DefaultGroup = "general"
	}

	result := make(GroupedEntries)

	for _, e := range entries {
		matched := false
		for _, prefix := range opts.Prefixes {
			if hasPrefix(e.Key, prefix) {
				result[prefix] = append(result[prefix], e)
				matched = true
				break
			}
		}
		if !matched {
			result[opts.DefaultGroup] = append(result[opts.DefaultGroup], e)
		}
	}

	return result
}

// GroupNames returns sorted group names from a GroupedEntries map.
func GroupNames(g GroupedEntries) []string {
	names := make([]string, 0, len(g))
	for k := range g {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// Flatten converts GroupedEntries back to a flat slice, preserving group order.
// Groups are iterated in sorted order.
func Flatten(g GroupedEntries) []Entry {
	var out []Entry
	for _, name := range GroupNames(g) {
		out = append(out, g[name]...)
	}
	return out
}

func hasPrefix(key, prefix string) bool {
	if len(key) <= len(prefix) {
		return false
	}
	return key[:len(prefix)+1] == prefix+"_"
}
