package envfile

import "fmt"

// DedupStrategy controls how duplicate keys are resolved.
type DedupStrategy int

const (
	// DedupKeepFirst retains the first occurrence of a duplicate key.
	DedupKeepFirst DedupStrategy = iota
	// DedupKeepLast retains the last occurrence of a duplicate key.
	DedupKeepLast
	// DedupError returns an error when a duplicate key is found.
	DedupError
)

// DedupResult holds the output of a Dedup operation.
type DedupResult struct {
	Entries  []Entry
	Removed  []Entry
	Duplicates map[string]int // key -> number of duplicates removed
}

// Dedup removes duplicate keys from entries according to the given strategy.
func Dedup(entries []Entry, strategy DedupStrategy) (*DedupResult, error) {
	seen := make(map[string]int) // key -> index in result
	result := &DedupResult{
		Duplicates: make(map[string]int),
	}

	for _, e := range entries {
		if idx, exists := seen[e.Key]; exists {
			switch strategy {
			case DedupError:
				return nil, fmt.Errorf("duplicate key: %q", e.Key)
			case DedupKeepLast:
				result.Removed = append(result.Removed, result.Entries[idx])
				result.Entries[idx] = e
				result.Duplicates[e.Key]++
			case DedupKeepFirst:
				result.Removed = append(result.Removed, e)
				result.Duplicates[e.Key]++
			}
		} else {
			seen[e.Key] = len(result.Entries)
			result.Entries = append(result.Entries, e)
		}
	}

	return result, nil
}

// HasDuplicates returns true if any key appears more than once.
func HasDuplicates(entries []Entry) bool {
	seen := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		if _, exists := seen[e.Key]; exists {
			return true
		}
		seen[e.Key] = struct{}{}
	}
	return false
}
