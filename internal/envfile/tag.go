package envfile

import (
	"fmt"
	"strings"
)

// Entry represents a key-value pair with optional metadata.
type TagOptions struct {
	// Keys to tag; if empty, all entries are tagged.
	Keys []string
	// Tags to add (e.g. "env:prod", "team:backend").
	Tags []string
	// Overwrite existing tags with the same name.
	Overwrite bool
}

// TagResult records what happened to each entry during tagging.
type TagResult struct {
	Key    string
	Tagged []string
	Skipped []string
}

// Tag adds metadata tags to matching entries via their comment field.
// Tags are stored as "# @tag:value" lines prepended to the entry comment.
func Tag(entries []Entry, opts TagOptions) ([]Entry, []TagResult, error) {
	if len(opts.Tags) == 0 {
		return nil, nil, fmt.Errorf("tag: at least one tag must be specified")
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	results := make([]TagResult, 0, len(entries))
	out := make([]Entry, len(entries))
	copy(out, entries)

	for i, e := range out {
		if len(opts.Keys) > 0 && !keySet[e.Key] {
			continue
		}
		res := TagResult{Key: e.Key}
		for _, tag := range opts.Tags {
			marker := "# @" + tag
			if strings.Contains(e.Comment, marker) && !opts.Overwrite {
				res.Skipped = append(res.Skipped, tag)
				continue
			}
			// Remove existing occurrence if overwriting.
			if opts.Overwrite {
				lines := strings.Split(e.Comment, "\n")
				filtered := lines[:0]
				for _, l := range lines {
					if !strings.HasPrefix(strings.TrimSpace(l), "# @"+tagName(tag)) {
						filtered = append(filtered, l)
					}
				}
				e.Comment = strings.Join(filtered, "\n")
			}
			if e.Comment == "" {
				e.Comment = marker
			} else {
				e.Comment = marker + "\n" + e.Comment
			}
			res.Tagged = append(res.Tagged, tag)
		}
		out[i] = e
		results = append(results, res)
	}
	return out, results, nil
}

// TaggedKeys returns keys that carry the given tag.
func TaggedKeys(entries []Entry, tag string) []string {
	marker := "# @" + tagName(tag)
	var keys []string
	for _, e := range entries {
		if strings.Contains(e.Comment, marker) {
			keys = append(keys, e.Key)
		}
	}
	return keys
}

// tagName extracts the tag name (before ':') for prefix matching.
func tagName(tag string) string {
	if idx := strings.Index(tag, ":"); idx >= 0 {
		return tag[:idx]
	}
	return tag
}
