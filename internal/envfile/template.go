package envfile

import (
	"fmt"
	"strings"
)

// GenerateTemplate produces a .env.template file from entries,
// replacing values with empty strings or placeholder comments.
// Sensitive keys (matched by the masker) get a hint comment.
func GenerateTemplate(entries []Entry, masker *Masker) []string {
	lines := make([]string, 0, len(entries)+2)
	lines = append(lines, "# Auto-generated template — fill in values before use")
	lines = append(lines, "")

	for _, e := range entries {
		if masker != nil && masker.IsSensitive(e.Key) {
			lines = append(lines, fmt.Sprintf("# %s — secret, do not commit real value", e.Key))
			lines = append(lines, fmt.Sprintf("%s=", e.Key))
		} else {
			lines = append(lines, fmt.Sprintf("%s=%s", e.Key, e.Value))
		}
	}
	return lines
}

// ApplyTemplate fills entries from a template using a values map.
// Keys present in values override the template defaults.
// Missing required keys (empty value in template) are returned as errors.
func ApplyTemplate(template []Entry, values map[string]string) ([]Entry, []error) {
	var result []Entry
	var errs []error

	for _, e := range template {
		val, ok := values[e.Key]
		if ok {
			result = append(result, Entry{Key: e.Key, Value: val})
			continue
		}
		if e.Value == "" {
			errs = append(errs, fmt.Errorf("required key %q has no value", e.Key))
			continue
		}
		result = append(result, e)
	}
	return result, errs
}

// TemplateKeys returns only the keys that have empty values (required placeholders).
func TemplateKeys(entries []Entry) []string {
	var keys []string
	for _, e := range entries {
		if strings.TrimSpace(e.Value) == "" {
			keys = append(keys, e.Key)
		}
	}
	return keys
}
