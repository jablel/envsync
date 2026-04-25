package envfile

import (
	"strings"
	"testing"
)

func TestLint_EmptyValue(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: ""}}
	result := Lint(entries)
	if !result.HasIssues() {
		t.Fatal("expected issues for empty value")
	}
	if result.Issues[0].Severity != LintWarn {
		t.Errorf("expected WARN, got %s", result.Issues[0].Severity)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	entries := []Entry{{Key: "my_key", Value: "val"}}
	result := Lint(entries)
	if !result.HasIssues() {
		t.Fatal("expected issues for lowercase key")
	}
	found := false
	for _, iss := range result.Issues {
		if strings.Contains(iss.Message, "uppercase") {
			found = true
		}
	}
	if !found {
		t.Error("expected uppercase warning")
	}
}

func TestLint_LeadingTrailingWhitespace(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: " bar "}}
	result := Lint(entries)
	if !result.HasIssues() {
		t.Fatal("expected issues for whitespace value")
	}
	if result.Issues[0].Severity != LintInfo {
		t.Errorf("expected INFO, got %s", result.Issues[0].Severity)
	}
}

func TestLint_Placeholder(t *testing.T) {
	entries := []Entry{{Key: "TOKEN", Value: "${{secrets.TOKEN}}"}}
	result := Lint(entries)
	if !result.HasIssues() {
		t.Fatal("expected issue for placeholder")
	}
	if !strings.Contains(result.Issues[0].Message, "placeholder") {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_CleanEntry(t *testing.T) {
	entries := []Entry{{Key: "DATABASE_URL", Value: "postgres://localhost/db"}}
	result := Lint(entries)
	if result.HasIssues() {
		t.Errorf("expected no issues, got: %v", result.Issues)
	}
}

func TestLint_LongValue(t *testing.T) {
	long := strings.Repeat("x", 600)
	entries := []Entry{{Key: "BIG_BLOB", Value: long}}
	result := Lint(entries)
	if !result.HasIssues() {
		t.Fatal("expected issue for long value")
	}
	if result.Issues[0].Severity != LintInfo {
		t.Errorf("expected INFO severity, got %s", result.Issues[0].Severity)
	}
}

func TestLintIssue_String(t *testing.T) {
	issue := LintIssue{Key: "FOO", Message: "value is empty", Severity: LintWarn}
	s := issue.String()
	if !strings.Contains(s, "WARN") || !strings.Contains(s, "FOO") {
		t.Errorf("unexpected string format: %s", s)
	}
}
