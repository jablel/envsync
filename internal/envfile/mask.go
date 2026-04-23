package envfile

import "strings"

// DefaultSecretPatterns contains common key substrings treated as secrets.
var DefaultSecretPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

const maskedValue = "***"

// Masker decides whether a key's value should be masked.
type Masker struct {
	Patterns []string
}

// NewMasker returns a Masker with the default secret patterns.
func NewMasker() *Masker {
	return &Masker{Patterns: DefaultSecretPatterns}
}

// IsSensitive returns true if the key matches any secret pattern.
func (m *Masker) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range m.Patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// MaskValue returns the masked placeholder if the key is sensitive,
// otherwise returns the original value unchanged.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return maskedValue
	}
	return value
}

// MaskEntries returns a copy of entries with sensitive values replaced.
func (m *Masker) MaskEntries(entries []Entry) []Entry {
	result := make([]Entry, len(entries))
	for i, e := range entries {
		result[i] = e
		if e.Key != "" {
			result[i].Value = m.MaskValue(e.Key, e.Value)
		}
	}
	return result
}
