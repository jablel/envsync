package envfile

import (
	"fmt"
)

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	OldKey  string
	NewKey  string
	Renamed bool
	Reason  string
}

// RenameOptions configures the behaviour of Rename.
type RenameOptions struct {
	// Overwrite allows the new key to overwrite an existing entry.
	Overwrite bool
}

// Rename renames oldKey to newKey in entries, returning the updated slice and
// a RenameResult describing what happened. The original slice is not mutated.
func Rename(entries []Entry, oldKey, newKey string, opts RenameOptions) ([]Entry, RenameResult, error) {
	if oldKey == "" {
		return entries, RenameResult{}, fmt.Errorf("oldKey must not be empty")
	}
	if newKey == "" {
		return entries, RenameResult{}, fmt.Errorf("newKey must not be empty")
	}
	if oldKey == newKey {
		return entries, RenameResult{OldKey: oldKey, NewKey: newKey, Renamed: false, Reason: "keys are identical"}, nil
	}

	oldIdx := -1
	newIdx := -1
	for i, e := range entries {
		if e.Key == oldKey {
			oldIdx = i
		}
		if e.Key == newKey {
			newIdx = i
		}
	}

	if oldIdx == -1 {
		return entries, RenameResult{OldKey: oldKey, NewKey: newKey, Renamed: false, Reason: "key not found"}, nil
	}

	if newIdx != -1 && !opts.Overwrite {
		return entries, RenameResult{}, fmt.Errorf("key %q already exists; use Overwrite to replace it", newKey)
	}

	// Build updated slice.
	updated := make([]Entry, 0, len(entries))
	for i, e := range entries {
		if i == newIdx {
			// Drop the existing target entry when overwriting.
			continue
		}
		if i == oldIdx {
			e.Key = newKey
		}
		updated = append(updated, e)
	}

	return updated, RenameResult{OldKey: oldKey, NewKey: newKey, Renamed: true}, nil
}
