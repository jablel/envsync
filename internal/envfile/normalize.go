package envfile

import (
	"strings"
)

// NormalizeOptions controls how entries are normalized.
type NormalizeOptions struct {
	// TrimWhitespace removes leading/trailing whitespace from values.
	TrimWhitespace bool
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// RemoveEmptyValues drops entries with empty values.
	RemoveEmptyValues bool
	// NormalizeLineEndings replaces \r\n and \r with \n in values.
	NormalizeLineEndings bool
}

// DefaultNormalizeOptions returns sensible defaults for normalization.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		TrimWhitespace:       true,
		UppercaseKeys:        false,
		RemoveEmptyValues:    false,
		NormalizeLineEndings: true,
	}
}

// NormalizeResult holds the outcome of a normalization pass.
type NormalizeResult struct {
	Entries  []Entry
	Changed  []string // keys that were modified
	Dropped  []string // keys that were removed
}

// Normalize applies the given options to a slice of entries, returning
// a NormalizeResult that describes what changed.
func Normalize(entries []Entry, opts NormalizeOptions) NormalizeResult {
	result := NormalizeResult{}

	for _, e := range entries {
		originalKey := e.Key
		originalValue := e.Value
		modified := false

		if opts.UppercaseKeys {
			upper := strings.ToUpper(e.Key)
			if upper != e.Key {
				e.Key = upper
				modified = true
			}
		}

		if opts.TrimWhitespace {
			trimmed := strings.TrimSpace(e.Value)
			if trimmed != e.Value {
				e.Value = trimmed
				modified = true
			}
		}

		if opts.NormalizeLineEndings {
			normalized := strings.ReplaceAll(e.Value, "\r\n", "\n")
			normalized = strings.ReplaceAll(normalized, "\r", "\n")
			if normalized != e.Value {
				e.Value = normalized
				modified = true
			}
		}

		if opts.RemoveEmptyValues && e.Value == "" {
			result.Dropped = append(result.Dropped, originalKey)
			continue
		}

		if modified {
			_ = originalValue
			result.Changed = append(result.Changed, e.Key)
		}

		result.Entries = append(result.Entries, e)
	}

	return result
}
