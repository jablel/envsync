package envfile

import (
	"fmt"
	"strings"
)

// LintSeverity represents the severity level of a lint warning.
type LintSeverity string

const (
	LintWarn  LintSeverity = "WARN"
	LintInfo  LintSeverity = "INFO"
)

// LintIssue describes a single linting finding.
type LintIssue struct {
	Key      string
	Message  string
	Severity LintSeverity
}

func (i LintIssue) String() string {
	return fmt.Sprintf("[%s] %s: %s", i.Severity, i.Key, i.Message)
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Issues []LintIssue
}

func (r *LintResult) HasIssues() bool {
	return len(r.Issues) > 0
}

func (r *LintResult) add(key, message string, severity LintSeverity) {
	r.Issues = append(r.Issues, LintIssue{Key: key, Message: message, Severity: severity})
}

// Lint checks entries for common style and correctness issues.
func Lint(entries []Entry) *LintResult {
	result := &LintResult{}

	for _, e := range entries {
		// Warn on empty values
		if e.Value == "" {
			result.add(e.Key, "value is empty", LintWarn)
		}

		// Warn on keys that are not UPPER_SNAKE_CASE
		if e.Key != strings.ToUpper(e.Key) {
			result.add(e.Key, "key is not uppercase", LintWarn)
		}

		// Info on values with leading/trailing whitespace (after unquoting)
		if strings.TrimSpace(e.Value) != e.Value {
			result.add(e.Key, "value has leading or trailing whitespace", LintInfo)
		}

		// Warn on values that look like they contain unexpanded placeholders
		if strings.Contains(e.Value, "${{") || strings.Contains(e.Value, "__PLACEHOLDER__") {
			result.add(e.Key, "value appears to contain an unexpanded placeholder", LintWarn)
		}

		// Info on very long values (possible accidental paste)
		if len(e.Value) > 512 {
			result.add(e.Key, "value is unusually long (>512 chars)", LintInfo)
		}
	}

	return result
}
