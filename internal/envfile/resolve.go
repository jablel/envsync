package envfile

import (
	"fmt"
	"strings"
)

// ResolveStrategy controls how conflicts are handled during resolution.
type ResolveStrategy int

const (
	// ResolvePreferBase keeps the base value on conflict.
	ResolvePreferBase ResolveStrategy = iota
	// ResolvePreferOverride uses the override value on conflict.
	ResolvePreferOverride
	// ResolveError returns an error on conflict.
	ResolveError
)

// ResolveOptions configures the Resolve operation.
type ResolveOptions struct {
	Strategy   ResolveStrategy
	EnvSuffix  string // optional: only apply overrides for keys ending with suffix
	StripSuffix bool   // strip the suffix from the resolved key name
}

// DefaultResolveOptions returns sensible defaults.
func DefaultResolveOptions() ResolveOptions {
	return ResolveOptions{
		Strategy:    ResolvePreferOverride,
		StripSuffix: false,
	}
}

// ResolveResult holds the outcome of a Resolve operation.
type ResolveResult struct {
	Entries   []Entry
	Conflicts []string
	Applied   int
}

// Resolve merges base and override entry slices according to options.
// Override entries take precedence unless the strategy says otherwise.
func Resolve(base, override []Entry, opts ResolveOptions) (ResolveResult, error) {
	baseMap := make(map[string]string, len(base))
	for _, e := range base {
		baseMap[e.Key] = e.Value
	}

	result := make([]Entry, 0, len(base))
	for _, e := range base {
		result = append(result, e)
	}

	indexMap := make(map[string]int, len(result))
	for i, e := range result {
		indexMap[e.Key] = i
	}

	var conflicts []string
	applied := 0

	for _, oe := range override {
		key := oe.Key
		if opts.EnvSuffix != "" {
			if !strings.HasSuffix(key, opts.EnvSuffix) {
				continue
			}
			if opts.StripSuffix {
				key = strings.TrimSuffix(key, opts.EnvSuffix)
			}
		}

		if idx, exists := indexMap[key]; exists {
			switch opts.Strategy {
			case ResolvePreferBase:
				conflicts = append(conflicts, key)
				continue
			case ResolveError:
				return ResolveResult{}, fmt.Errorf("resolve conflict on key %q", key)
			default:
				result[idx] = Entry{Key: key, Value: oe.Value}
				applied++
			}
		} else {
			result = append(result, Entry{Key: key, Value: oe.Value})
			indexMap[key] = len(result) - 1
			applied++
		}
	}

	return ResolveResult{Entries: result, Conflicts: conflicts, Applied: applied}, nil
}
