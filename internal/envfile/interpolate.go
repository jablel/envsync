package envfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// interpolatePattern matches ${VAR} and $VAR style references.
var interpolatePattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateOptions controls how variable interpolation behaves.
type InterpolateOptions struct {
	// AllowEnvFallback lets unresolved variables fall back to OS environment.
	AllowEnvFallback bool
	// FailOnMissing returns an error if a referenced variable cannot be resolved.
	FailOnMissing bool
}

// DefaultInterpolateOptions returns sensible defaults.
func DefaultInterpolateOptions() InterpolateOptions {
	return InterpolateOptions{
		AllowEnvFallback: true,
		FailOnMissing:    false,
	}
}

// Interpolate resolves variable references within entry values using the
// provided entries as the source of truth, with optional OS env fallback.
// It returns a new slice of entries with values expanded.
func Interpolate(entries []Entry, opts InterpolateOptions) ([]Entry, error) {
	// Build a lookup map from the provided entries.
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}

	result := make([]Entry, len(entries))
	for i, e := range entries {
		expanded, err := interpolateValue(e.Value, lookup, opts)
		if err != nil {
			return nil, fmt.Errorf("interpolate: key %q: %w", e.Key, err)
		}
		result[i] = Entry{Key: e.Key, Value: expanded, Comment: e.Comment}
	}
	return result, nil
}

// interpolateValue expands variable references within a single value string.
func interpolateValue(value string, lookup map[string]string, opts InterpolateOptions) (string, error) {
	var resolveErr error

	expanded := interpolatePattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}

		// Extract variable name from either ${VAR} or $VAR form.
		sub := interpolatePattern.FindStringSubmatch(match)
		varName := sub[1]
		if varName == "" {
			varName = sub[2]
		}

		if v, ok := lookup[varName]; ok {
			return v
		}
		if opts.AllowEnvFallback {
			if v, ok := os.LookupEnv(varName); ok {
				return v
			}
		}
		if opts.FailOnMissing {
			resolveErr = fmt.Errorf("unresolved variable %q", varName)
			return match
		}
		return strings.TrimSpace(match)
	})

	if resolveErr != nil {
		return "", resolveErr
	}
	return expanded, nil
}
