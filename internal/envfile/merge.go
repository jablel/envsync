package envfile

import "fmt"

// MergeStrategy defines how conflicting keys are handled during merge.
type MergeStrategy int

const (
	// MergeStrategyKeepBase keeps the base value on conflict.
	MergeStrategyKeepBase MergeStrategy = iota
	// MergeStrategyKeepOverride uses the override value on conflict.
	MergeStrategyKeepOverride
	// MergeStrategyError returns an error on conflict.
	MergeStrategyError
)

// MergeOptions configures the merge behaviour.
type MergeOptions struct {
	Strategy      MergeStrategy
	SkipBlanks    bool
	SkipComments  bool
}

// MergeResult holds the merged entries and a list of conflicts that were
// resolved according to the chosen strategy.
type MergeResult struct {
	Entries   []Entry
	Conflicts []string // keys that had conflicting values
}

// Merge combines base and override env entries according to opts.
// The order of keys in base is preserved; new keys from override are
// appended in the order they appear.
func Merge(base, override []Entry, opts MergeOptions) (*MergeResult, error) {
	baseMap := make(map[string]int, len(base)) // key -> index in result
	result := make([]Entry, 0, len(base))

	for _, e := range base {
		if opts.SkipBlanks && e.Key == "" && e.Value == "" {
			continue
		}
		baseMap[e.Key] = len(result)
		result = append(result, e)
	}

	var conflicts []string

	for _, e := range override {
		if opts.SkipBlanks && e.Key == "" && e.Value == "" {
			continue
		}
		if e.Key == "" {
			// comment or blank line — append only if not skipping
			if !opts.SkipComments {
				result = append(result, e)
			}
			continue
		}

		idx, exists := baseMap[e.Key]
		if !exists {
			baseMap[e.Key] = len(result)
			result = append(result, e)
			continue
		}

		// Conflict: key exists in both
		if result[idx].Value == e.Value {
			continue // same value, no real conflict
		}

		switch opts.Strategy {
		case MergeStrategyKeepBase:
			conflicts = append(conflicts, e.Key)
		case MergeStrategyKeepOverride:
			conflicts = append(conflicts, e.Key)
			result[idx].Value = e.Value
		case MergeStrategyError:
			return nil, fmt.Errorf("merge conflict on key %q: base=%q override=%q", e.Key, result[idx].Value, e.Value)
		}
	}

	return &MergeResult{Entries: result, Conflicts: conflicts}, nil
}
