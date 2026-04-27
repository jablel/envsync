package envfile

import (
	"fmt"
	"strings"
)

// RotateSummary returns a human-readable summary of a RotateResult.
func RotateSummary(result RotateResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Rotated: %d key(s)\n", len(result.Rotated)))
	for _, k := range result.Rotated {
		sb.WriteString(fmt.Sprintf("  ~ %s\n", k))
	}

	if len(result.Errors) > 0 {
		sb.WriteString(fmt.Sprintf("Failed:  %d key(s)\n", len(result.Errors)))
		for k, e := range result.Errors {
			sb.WriteString(fmt.Sprintf("  ! %s: %v\n", k, e))
		}
	}

	sb.WriteString(fmt.Sprintf("Skipped: %d key(s)\n", len(result.Skipped)))
	return strings.TrimRight(sb.String(), "\n")
}

// FormatRotateInstructions returns a patch-style instruction list for the rotation.
func FormatRotateInstructions(original, updated []Entry, masker *Masker) string {
	var sb strings.Builder
	orig := make(map[string]string, len(original))
	for _, e := range original {
		orig[e.Key] = e.Value
	}

	for _, e := range updated {
		old, ok := orig[e.Key]
		if !ok || old == e.Value {
			continue
		}
		oldDisplay := masker.MaskValue(e.Key, old)
		newDisplay := masker.MaskValue(e.Key, e.Value)
		sb.WriteString(fmt.Sprintf("rotate %s: %s -> %s\n", e.Key, oldDisplay, newDisplay))
	}
	return strings.TrimRight(sb.String(), "\n")
}
