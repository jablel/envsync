package envfile

import (
	"fmt"
)

// CloneOptions controls how entries are cloned between environments.
type CloneOptions struct {
	// KeyFilter, if non-nil, is called for each entry; return true to include it.
	KeyFilter func(key string) bool
	// Overwrite controls whether existing keys in dst are replaced.
	Overwrite bool
	// StripSecrets replaces sensitive values with empty strings in the clone.
	StripSecrets bool
}

// CloneResult summarises what happened during a Clone operation.
type CloneResult struct {
	Copied   []string
	Skipped  []string
	Stripped []string
}

// Clone copies entries from src into dst according to opts.
// It returns a new slice (dst is not mutated) and a CloneResult summary.
func Clone(src, dst []Entry, opts CloneOptions) ([]Entry, CloneResult, error) {
	if src == nil {
		return nil, CloneResult{}, fmt.Errorf("clone: src must not be nil")
	}

	masker := NewMasker(nil)

	// Build a lookup map for dst so we can detect conflicts quickly.
	dstMap := make(map[string]int, len(dst))
	for i, e := range dst {
		dstMap[e.Key] = i
	}

	out := make([]Entry, len(dst))
	copy(out, dst)

	var result CloneResult

	for _, e := range src {
		if opts.KeyFilter != nil && !opts.KeyFilter(e.Key) {
			continue
		}

		entry := e

		if opts.StripSecrets && masker.IsSensitive(e.Key) {
			entry.Value = ""
			result.Stripped = append(result.Stripped, e.Key)
		}

		if idx, exists := dstMap[e.Key]; exists {
			if !opts.Overwrite {
				result.Skipped = append(result.Skipped, e.Key)
				continue
			}
			out[idx] = entry
		} else {
			dstMap[e.Key] = len(out)
			out = append(out, entry)
		}

		result.Copied = append(result.Copied, e.Key)
	}

	return out, result, nil
}
