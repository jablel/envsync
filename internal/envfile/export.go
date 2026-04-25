package envfile

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ExportFormat defines the output format for exporting env entries.
type ExportFormat string

const (
	FormatDotEnv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
	FormatShell  ExportFormat = "shell"
)

// ExportOptions controls how entries are exported.
type ExportOptions struct {
	Format    ExportFormat
	SortKeys  bool
	MaskSecrets bool
	Masker    *Masker
}

// Export converts a slice of Entry into the requested format string.
func Export(entries []Entry, opts ExportOptions) (string, error) {
	if opts.MaskSecrets && opts.Masker == nil {
		opts.Masker = NewMasker(nil)
	}

	working := make([]Entry, len(entries))
	copy(working, entries)

	if opts.SortKeys {
		sort.Slice(working, func(i, j int) bool {
			return working[i].Key < working[j].Key
		})
	}

	if opts.MaskSecrets && opts.Masker != nil {
		working = opts.Masker.MaskEntries(working)
	}

	switch opts.Format {
	case FormatJSON:
		return exportJSON(working)
	case FormatShell:
		return exportShell(working)
	case FormatDotEnv, "":
		return exportDotEnv(working), nil
	default:
		return "", fmt.Errorf("unsupported export format: %q", opts.Format)
	}
}

func exportDotEnv(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		if needsQuoting(e.Value) {
			fmt.Fprintf(&sb, "%s=%q\n", e.Key, e.Value)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
		}
	}
	return sb.String()
}

func exportJSON(entries []Entry) (string, error) {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(b) + "\n", nil
}

func exportShell(entries []Entry) (string, error) {
	var sb strings.Builder
	for _, e := range entries {
		if needsQuoting(e.Value) {
			fmt.Fprintf(&sb, "export %s=%q\n", e.Key, e.Value)
		} else {
			fmt.Fprintf(&sb, "export %s=%s\n", e.Key, e.Value)
		}
	}
	return sb.String(), nil
}
