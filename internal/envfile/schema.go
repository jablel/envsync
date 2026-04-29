package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaField defines the expected shape of a single env entry.
type SchemaField struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern the value must match
	Default  string // used in reporting only
}

// SchemaResult holds the outcome of validating entries against a schema.
type SchemaResult struct {
	Missing  []string // required keys absent from entries
	Invalid  []string // keys present but value fails pattern
	Extra    []string // keys in entries not defined in schema
	Valid    []string // keys that passed all checks
}

// HasErrors returns true if there are any missing or invalid fields.
func (r SchemaResult) HasErrors() bool {
	return len(r.Missing) > 0 || len(r.Invalid) > 0
}

// Summary returns a human-readable summary of the schema validation result.
func (r SchemaResult) Summary() string {
	var sb strings.Builder
	if len(r.Missing) > 0 {
		fmt.Fprintf(&sb, "missing required keys: %s\n", strings.Join(r.Missing, ", "))
	}
	if len(r.Invalid) > 0 {
		fmt.Fprintf(&sb, "invalid values for keys: %s\n", strings.Join(r.Invalid, ", "))
	}
	if len(r.Extra) > 0 {
		fmt.Fprintf(&sb, "extra keys not in schema: %s\n", strings.Join(r.Extra, ", "))
	}
	fmt.Fprintf(&sb, "valid: %d key(s)\n", len(r.Valid))
	return strings.TrimRight(sb.String(), "\n")
}

// ValidateSchema checks a slice of Entry values against the provided schema fields.
func ValidateSchema(entries []Entry, fields []SchemaField) (SchemaResult, error) {
	result := SchemaResult{}

	entryMap := make(map[string]string, len(entries))
	for _, e := range entries {
		entryMap[e.Key] = e.Value
	}

	schemaKeys := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		schemaKeys[f.Key] = struct{}{}
		val, exists := entryMap[f.Key]
		if !exists {
			if f.Required {
				result.Missing = append(result.Missing, f.Key)
			}
			continue
		}
		if f.Pattern != "" {
			re, err := regexp.Compile(f.Pattern)
			if err != nil {
				return result, fmt.Errorf("invalid pattern for key %q: %w", f.Key, err)
			}
			if !re.MatchString(val) {
				result.Invalid = append(result.Invalid, f.Key)
				continue
			}
		}
		result.Valid = append(result.Valid, f.Key)
	}

	for _, e := range entries {
		if _, inSchema := schemaKeys[e.Key]; !inSchema {
			result.Extra = append(result.Extra, e.Key)
		}
	}

	return result, nil
}
