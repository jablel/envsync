package envfile

import (
	"fmt"
	"strings"
)

// RotateOptions controls how key rotation is performed.
type RotateOptions struct {
	// Keys is the explicit list of keys to rotate. If empty, all sensitive keys are rotated.
	Keys []string
	// Generator is called with the key name and returns a new value.
	Generator func(key string) (string, error)
	// DryRun reports what would change without modifying entries.
	DryRun bool
}

// RotateResult holds the outcome of a rotation operation.
type RotateResult struct {
	Rotated []string
	Skipped []string
	Errors  map[string]error
}

// Rotate replaces values for sensitive (or specified) keys using the provided generator.
// It returns the updated entries and a RotateResult summarising changes.
func Rotate(entries []Entry, masker *Masker, opts RotateOptions) ([]Entry, RotateResult, error) {
	if opts.Generator == nil {
		return nil, RotateResult{}, fmt.Errorf("rotate: Generator must not be nil")
	}

	targetSet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		targetSet[k] = true
	}

	result := RotateResult{
		Errors: make(map[string]error),
	}

	updated := make([]Entry, len(entries))
	copy(updated, entries)

	for i, e := range updated {
		shouldRotate := len(targetSet) > 0 && targetSet[e.Key]
		if len(targetSet) == 0 {
			shouldRotate = masker.IsSensitive(e.Key)
		}

		if !shouldRotate {
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}

		newVal, err := opts.Generator(e.Key)
		if err != nil {
			result.Errors[e.Key] = err
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}

		if !opts.DryRun {
			updated[i].Value = strings.TrimSpace(newVal)
		}
		result.Rotated = append(result.Rotated, e.Key)
	}

	if len(result.Errors) > 0 {
		return updated, result, fmt.Errorf("rotate: %d key(s) failed to generate new values", len(result.Errors))
	}
	return updated, result, nil
}
