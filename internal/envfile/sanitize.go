package envfile

import (
	"fmt"
	"strings"
)

// SanitizeOptions controls which sanitization rules are applied.
type SanitizeOptions struct {
	StripControlChars bool // remove non-printable control characters from values
	TrimQuotes        bool // remove surrounding quotes left in values
	NormalizeKeys     bool // uppercase and replace invalid chars with underscores
	RemoveNullBytes   bool // strip null bytes from values
}

// DefaultSanitizeOptions returns sensible defaults for sanitization.
func DefaultSanitizeOptions() SanitizeOptions {
	return SanitizeOptions{
		StripControlChars: true,
		TrimQuotes:        false,
		NormalizeKeys:     false,
		RemoveNullBytes:   true,
	}
}

// SanitizeResult holds the outcome of a sanitize operation on a single entry.
type SanitizeResult struct {
	Key     string
	Before  string
	After   string
	Changed bool
}

// Sanitize cleans up env entries according to the provided options.
// It returns the sanitized entries and a per-entry result slice.
func Sanitize(entries []Entry, opts SanitizeOptions) ([]Entry, []SanitizeResult, error) {
	results := make([]SanitizeResult, 0, len(entries))
	out := make([]Entry, 0, len(entries))

	for _, e := range entries {
		origKey := e.Key
		origVal := e.Value

		if opts.RemoveNullBytes {
			e.Value = strings.ReplaceAll(e.Value, "\x00", "")
		}

		if opts.StripControlChars {
			e.Value = stripControlChars(e.Value)
		}

		if opts.TrimQuotes {
			e.Value = strings.Trim(e.Value, `"'`)
		}

		if opts.NormalizeKeys {
			newKey, err := normalizeKey(e.Key)
			if err != nil {
				return nil, nil, fmt.Errorf("sanitize: key %q: %w", e.Key, err)
			}
			e.Key = newKey
		}

		changed := e.Key != origKey || e.Value != origVal
		results = append(results, SanitizeResult{
			Key:     origKey,
			Before:  origVal,
			After:   e.Value,
			Changed: changed,
		})
		out = append(out, e)
	}

	return out, results, nil
}

// SanitizedKeys returns the keys of entries that were modified.
func SanitizedKeys(results []SanitizeResult) []string {
	var keys []string
	for _, r := range results {
		if r.Changed {
			keys = append(keys, r.Key)
		}
	}
	return keys
}

func stripControlChars(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 0x20 || r == '\t' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func normalizeKey(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("empty key")
	}
	key = strings.ToUpper(key)
	var b strings.Builder
	for i, r := range key {
		switch {
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9' && i > 0:
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String(), nil
}
