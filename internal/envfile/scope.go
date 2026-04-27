package envfile

import "fmt"

// ScopeOptions controls how scoping is applied.
type ScopeOptions struct {
	// Prefix is prepended to all keys (e.g. "APP_").
	Prefix string
	// Strip removes the prefix from all keys instead of adding it.
	Strip bool
	// SkipNoMatch silently skips keys that don't match the prefix when stripping.
	SkipNoMatch bool
}

// ScopeResult holds the output of a Scope operation.
type ScopeResult struct {
	Entries  []Entry
	Skipped  []Entry
}

// Scope applies or removes a namespace prefix from env entries.
//
// When Strip is false, each key is prefixed with Prefix.
// When Strip is true, Prefix is removed from matching keys; non-matching
// keys are either skipped (SkipNoMatch=true) or returned as-is.
func Scope(entries []Entry, opts ScopeOptions) (ScopeResult, error) {
	if opts.Prefix == "" {
		return ScopeResult{}, fmt.Errorf("scope: prefix must not be empty")
	}

	result := ScopeResult{}

	for _, e := range entries {
		if opts.Strip {
			if len(e.Key) >= len(opts.Prefix) && e.Key[:len(opts.Prefix)] == opts.Prefix {
				scoped := e
				scoped.Key = e.Key[len(opts.Prefix):]
				result.Entries = append(result.Entries, scoped)
			} else {
				if opts.SkipNoMatch {
					result.Skipped = append(result.Skipped, e)
				} else {
					result.Entries = append(result.Entries, e)
				}
			}
		} else {
			scoped := e
			scoped.Key = opts.Prefix + e.Key
			result.Entries = append(result.Entries, scoped)
		}
	}

	return result, nil
}

// ScopeKeys returns only the keys that would result from a Scope operation.
func ScopeKeys(entries []Entry, opts ScopeOptions) ([]string, error) {
	res, err := Scope(entries, opts)
	if err != nil {
		return nil, err
	}
	keys := make([]string, len(res.Entries))
	for i, e := range res.Entries {
		keys[i] = e.Key
	}
	return keys, nil
}
