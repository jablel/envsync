package envfile

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a single Entry.
type TransformFunc func(Entry) (Entry, error)

// TransformOptions controls the behaviour of Transform.
type TransformOptions struct {
	// SkipErrors causes entries that produce a transform error to be left
	// unchanged rather than aborting the whole operation.
	SkipErrors bool
}

// Transform applies one or more TransformFuncs to every entry in order.
// If multiple funcs are provided they are chained: the output of one becomes
// the input of the next.
func Transform(entries []Entry, opts TransformOptions, fns ...TransformFunc) ([]Entry, error) {
	if len(fns) == 0 {
		return entries, nil
	}

	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		current := e
		var err error
		for _, fn := range fns {
			current, err = fn(current)
			if err != nil {
				if opts.SkipErrors {
					current = e // restore original
					break
				}
				return nil, fmt.Errorf("transform error on key %q: %w", e.Key, err)
			}
		}
		result = append(result, current)
	}
	return result, nil
}

// UppercaseKeys returns a TransformFunc that converts every key to upper-case.
func UppercaseKeys() TransformFunc {
	return func(e Entry) (Entry, error) {
		e.Key = strings.ToUpper(e.Key)
		return e, nil
	}
}

// PrefixKeys returns a TransformFunc that prepends prefix to every key.
func PrefixKeys(prefix string) TransformFunc {
	return func(e Entry) (Entry, error) {
		if prefix == "" {
			return e, fmt.Errorf("prefix must not be empty")
		}
		e.Key = prefix + e.Key
		return e, nil
	}
}

// TrimValues returns a TransformFunc that strips leading/trailing whitespace
// from every value.
func TrimValues() TransformFunc {
	return func(e Entry) (Entry, error) {
		e.Value = strings.TrimSpace(e.Value)
		return e, nil
	}
}

// MaskTransform returns a TransformFunc that masks sensitive values using the
// provided Masker.
func MaskTransform(m *Masker) TransformFunc {
	return func(e Entry) (Entry, error) {
		if m.IsSensitive(e.Key) {
			e.Value = m.MaskValue(e.Value)
		}
		return e, nil
	}
}
