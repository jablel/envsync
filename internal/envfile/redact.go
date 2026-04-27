package envfile

import (
	"strings"
)

// RedactOptions controls how redaction is applied.
type RedactOptions struct {
	// Replacement is the string used to replace sensitive values. Defaults to "[REDACTED]".
	Replacement string
	// PartialReveal reveals the first N characters of a sensitive value before redacting the rest.
	PartialReveal int
}

// DefaultRedactOptions returns sensible defaults for RedactOptions.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Replacement:   "[REDACTED]",
		PartialReveal: 0,
	}
}

// Redact returns a copy of entries where sensitive values are replaced
// according to the provided options and masker.
func Redact(entries []Entry, masker *Masker, opts RedactOptions) []Entry {
	if opts.Replacement == "" {
		opts.Replacement = "[REDACTED]"
	}

	result := make([]Entry, len(entries))
	for i, e := range entries {
		if masker.IsSensitive(e.Key) {
			e.Value = redactValue(e.Value, opts)
		}
		result[i] = e
	}
	return result
}

// RedactMap returns a copy of a key→value map with sensitive values redacted.
func RedactMap(m map[string]string, masker *Masker, opts RedactOptions) map[string]string {
	if opts.Replacement == "" {
		opts.Replacement = "[REDACTED]"
	}

	out := make(map[string]string, len(m))
	for k, v := range m {
		if masker.IsSensitive(k) {
			out[k] = redactValue(v, opts)
		} else {
			out[k] = v
		}
	}
	return out
}

// redactValue applies partial reveal logic then replaces the remainder.
func redactValue(value string, opts RedactOptions) string {
	if opts.PartialReveal <= 0 || len(value) <= opts.PartialReveal {
		return opts.Replacement
	}
	visible := value[:opts.PartialReveal]
	return visible + strings.Repeat("*", len(value)-opts.PartialReveal)
}
