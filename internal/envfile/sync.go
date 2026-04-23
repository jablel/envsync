package envfile

import (
	"fmt"
	"os"
	"strings"
)

// SyncResult holds the outcome of a sync operation.
type SyncResult struct {
	Applied  []string
	Skipped  []string
	Errors   []string
}

// SyncOptions controls the behavior of the Sync function.
type SyncOptions struct {
	// Overwrite existing keys in the target if true.
	Overwrite bool
	// DryRun reports what would change without writing.
	DryRun bool
}

// Sync merges entries from source into target according to opts.
// It writes the updated target file to targetPath unless DryRun is set.
func Sync(source, target map[string]string, targetPath string, opts SyncOptions) (*SyncResult, error) {
	result := &SyncResult{}

	updated := make(map[string]string, len(target))
	for k, v := range target {
		updated[k] = v
	}

	for key, val := range source {
		if _, exists := updated[key]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, key)
			continue
		}
		updated[key] = val
		result.Applied = append(result.Applied, key)
	}

	if opts.DryRun {
		return result, nil
	}

	if err := writeEnvFile(targetPath, updated); err != nil {
		return result, fmt.Errorf("sync: write target: %w", err)
	}

	return result, nil
}

// writeEnvFile serializes entries to a .env formatted file.
func writeEnvFile(path string, entries map[string]string) error {
	var sb strings.Builder
	for k, v := range entries {
		if strings.ContainsAny(v, " \t\n#") {
			fmt.Fprintf(&sb, "%s=\"%s\"\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return os.WriteFile(path, []byte(sb.String()), 0o600)
}
