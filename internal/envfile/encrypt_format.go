package envfile

import (
	"fmt"
	"strings"
)

// EncryptSummary returns a human-readable summary of an EncryptResult.
func EncryptSummary(res EncryptResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Encrypted: %d key(s), Skipped: %d key(s)\n",
		len(res.Encrypted), len(res.Skipped)))
	if len(res.Encrypted) > 0 {
		sb.WriteString("  Encrypted keys:\n")
		for _, e := range res.Encrypted {
			sb.WriteString(fmt.Sprintf("    + %s\n", e.Key))
		}
	}
	if len(res.Skipped) > 0 {
		sb.WriteString("  Skipped keys:\n")
		for _, e := range res.Skipped {
			sb.WriteString(fmt.Sprintf("    - %s\n", e.Key))
		}
	}
	return sb.String()
}

// FormatEncryptedEntries formats entries, masking encrypted values.
func FormatEncryptedEntries(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		val := e.Value
		if strings.HasPrefix(val, encPrefix) {
			val = encPrefix + "<encrypted>"
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, val))
	}
	return sb.String()
}
