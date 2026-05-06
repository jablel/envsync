package envfile

// CompactOptions controls how Compact removes redundant entries.
type CompactOptions struct {
	// RemoveCommented removes entries whose values are commented-out placeholders.
	RemoveCommented bool
	// RemoveDuplicates keeps only the last occurrence of each key.
	RemoveDuplicates bool
	// RemoveEmpty removes entries with empty values.
	RemoveEmpty bool
	// RemoveDefaults removes entries whose value matches the provided defaults map.
	RemoveDefaults map[string]string
}

// CompactResult describes what Compact did.
type CompactResult struct {
	Output   []Entry
	Removed  []Entry
	Kept     int
}

// DefaultCompactOptions returns a CompactOptions with safe defaults.
func DefaultCompactOptions() CompactOptions {
	return CompactOptions{
		RemoveDuplicates: true,
	}
}

// Compact reduces an env entry slice by removing redundant or unwanted entries
// according to the provided options.
func Compact(entries []Entry, opts CompactOptions) CompactResult {
	var removed []Entry
	working := make([]Entry, 0, len(entries))

	// Pass 1: remove empty and commented values.
	for _, e := range entries {
		if opts.RemoveEmpty && e.Value == "" {
			removed = append(removed, e)
			continue
		}
		if opts.RemoveCommented && isCommentedPlaceholder(e.Value) {
			removed = append(removed, e)
			continue
		}
		if opts.RemoveDefaults != nil {
			if def, ok := opts.RemoveDefaults[e.Key]; ok && e.Value == def {
				removed = append(removed, e)
				continue
			}
		}
		working = append(working, e)
	}

	// Pass 2: deduplicate, keeping last occurrence.
	if opts.RemoveDuplicates {
		seen := make(map[string]int, len(working))
		for i, e := range working {
			seen[e.Key] = i
		}
		deduped := make([]Entry, 0, len(working))
		visited := make(map[string]bool)
		for i := len(working) - 1; i >= 0; i-- {
			e := working[i]
			if visited[e.Key] {
				removed = append(removed, e)
				continue
			}
			if seen[e.Key] == i {
				deduped = append([]Entry{e}, deduped...)
				visited[e.Key] = true
			}
		}
		working = deduped
	}

	return CompactResult{
		Output:  working,
		Removed: removed,
		Kept:    len(working),
	}
}

// isCommentedPlaceholder returns true if the value looks like a commented-out
// placeholder such as "#changeme" or "# TODO".
func isCommentedPlaceholder(v string) bool {
	if len(v) == 0 {
		return false
	}
	return v[0] == '#'
}
