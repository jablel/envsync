package envfile

import (
	"fmt"
	"sort"
)

// FreezeOptions controls which keys are frozen and how conflicts are handled.
type FreezeOptions struct {
	// Keys is an explicit list of keys to freeze. If empty, all keys are frozen.
	Keys []string
	// AllowUnfreeze permits overwriting a frozen value when set to true.
	AllowUnfreeze bool
}

// FreezeResult captures the outcome of a Freeze operation.
type FreezeResult struct {
	Frozen  []string
	Skipped []string
	Errors  []string
}

// frozenTag is the comment prefix used to mark frozen entries.
const frozenTag = "@frozen"

// Freeze marks the specified entries as frozen by appending a @frozen tag to
// their comment. Subsequent calls to Freeze (or other write operations that
// respect frozen state) will refuse to modify those entries unless
// AllowUnfreeze is true.
func Freeze(entries []Entry, opts FreezeOptions) ([]Entry, FreezeResult, error) {
	target := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		target[k] = true
	}

	result := FreezeResult{}
	out := make([]Entry, len(entries))

	for i, e := range entries {
		if len(opts.Keys) > 0 && !target[e.Key] {
			out[i] = e
			continue
		}

		if IsFrozen(e) {
			if !opts.AllowUnfreeze {
				result.Skipped = append(result.Skipped, e.Key)
				out[i] = e
				continue
			}
			// re-freeze is a no-op value-wise; just record it
		}

		e.Comment = setFrozenTag(e.Comment)
		out[i] = e
		result.Frozen = append(result.Frozen, e.Key)
	}

	sort.Strings(result.Frozen)
	sort.Strings(result.Skipped)
	return out, result, nil
}

// IsFrozen reports whether the given entry carries the @frozen tag.
func IsFrozen(e Entry) bool {
	return containsTag(e.Comment, frozenTag)
}

// FrozenKeys returns the keys of all frozen entries.
func FrozenKeys(entries []Entry) []string {
	var keys []string
	for _, e := range entries {
		if IsFrozen(e) {
			keys = append(keys, e.Key)
		}
	}
	return keys
}

func setFrozenTag(comment string) string {
	if containsTag(comment, frozenTag) {
		return comment
	}
	if comment == "" {
		return fmt.Sprintf("# %s", frozenTag)
	}
	return fmt.Sprintf("%s %s", comment, frozenTag)
}

func containsTag(comment, tag string) bool {
	return len(comment) >= len(tag) && (comment == "# "+tag ||
		len(comment) > len(tag)+1 && comment[len(comment)-len(tag):] == tag)
}
