package envfile

import "fmt"

// PromoteOptions controls how keys are promoted between environments.
type PromoteOptions struct {
	// Overwrite replaces existing keys in the target.
	Overwrite bool
	// DryRun reports what would change without writing.
	DryRun bool
	// Keys restricts promotion to only these keys. If empty, all keys are promoted.
	Keys []string
}

// PromoteResult describes the outcome of a promotion operation.
type PromoteResult struct {
	Promoted []string
	Skipped  []string
	Overwritten []string
}

// Promote copies selected (or all) keys from src into dst according to opts.
// It returns the updated dst entries and a PromoteResult summary.
func Promote(src, dst []Entry, opts PromoteOptions) ([]Entry, PromoteResult, error) {
	filter := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		filter[k] = true
	}

	dstMap := make(map[string]int, len(dst))
	for i, e := range dst {
		dstMap[e.Key] = i
	}

	result := PromoteResult{}
	updated := make([]Entry, len(dst))
	copy(updated, dst)

	for _, e := range src {
		if len(filter) > 0 && !filter[e.Key] {
			continue
		}
		if idx, exists := dstMap[e.Key]; exists {
			if !opts.Overwrite {
				result.Skipped = append(result.Skipped, e.Key)
				continue
			}
			if !opts.DryRun {
				updated[idx] = e
			}
			result.Overwritten = append(result.Overwritten, e.Key)
		} else {
			if !opts.DryRun {
				updated = append(updated, e)
				dstMap[e.Key] = len(updated) - 1
			}
			result.Promoted = append(result.Promoted, e.Key)
		}
	}

	if len(filter) > 0 {
		for k := range filter {
			found := false
			for _, e := range src {
				if e.Key == k {
					found = true
					break
				}
			}
			if !found {
				return nil, PromoteResult{}, fmt.Errorf("promote: key %q not found in source", k)
			}
		}
	}

	return updated, result, nil
}
