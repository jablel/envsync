package envfile

import "fmt"

// ReorderOptions controls how entries are reordered.
type ReorderOptions struct {
	// Keys defines the desired key order. Keys not listed appear at the end.
	Keys []string
	// UnlistedAtEnd places unlisted keys after listed ones (default true).
	UnlistedAtEnd bool
}

// DefaultReorderOptions returns sensible defaults.
func DefaultReorderOptions() ReorderOptions {
	return ReorderOptions{
		UnlistedAtEnd: true,
	}
}

// ReorderResult holds the outcome of a Reorder call.
type ReorderResult struct {
	Entries  []Entry
	Moved    []string // keys whose position changed
	Missing  []string // keys listed in order but absent from entries
}

// Reorder rearranges entries according to the provided key order.
// Keys not present in opts.Keys are appended at the end when UnlistedAtEnd
// is true, otherwise they are dropped.
func Reorder(entries []Entry, opts ReorderOptions) (ReorderResult, error) {
	if len(opts.Keys) == 0 {
		return ReorderResult{}, fmt.Errorf("reorder: no key order specified")
	}

	index := make(map[string]Entry, len(entries))
	origPos := make(map[string]int, len(entries))
	for i, e := range entries {
		index[e.Key] = e
		origPos[e.Key] = i
	}

	seen := make(map[string]bool)
	result := make([]Entry, 0, len(entries))
	var missing []string

	for _, k := range opts.Keys {
		if e, ok := index[k]; ok {
			result = append(result, e)
			seen[k] = true
		} else {
			missing = append(missing, k)
		}
	}

	if opts.UnlistedAtEnd {
		for _, e := range entries {
			if !seen[e.Key] {
				result = append(result, e)
			}
		}
	}

	var moved []string
	for newIdx, e := range result {
		if origPos[e.Key] != newIdx {
			moved = append(moved, e.Key)
		}
	}

	return ReorderResult{
		Entries: result,
		Moved:   moved,
		Missing: missing,
	}, nil
}
