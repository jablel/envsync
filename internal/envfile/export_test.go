package envfile

import (
	"strings"
	"testing"
)

func TestExport_DotEnv(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "PORT", Value: "8080"},
	}
	out, err := Export(entries, ExportOptions{Format: FormatDotEnv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME=myapp in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got:\n%s", out)
	}
}

func TestExport_JSON(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DEBUG", Value: "true"},
	}
	out, err := Export(entries, ExportOptions{Format: FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"APP_NAME"`) {
		t.Errorf("expected APP_NAME key in JSON, got:\n%s", out)
	}
	if !strings.Contains(out, `"myapp"`) {
		t.Errorf("expected myapp value in JSON, got:\n%s", out)
	}
}

func TestExport_Shell(t *testing.T) {
	entries := []Entry{
		{Key: "PATH_EXTRA", Value: "/usr/local/bin"},
	}
	out, err := Export(entries, ExportOptions{Format: FormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export PATH_EXTRA=") {
		t.Errorf("expected 'export PATH_EXTRA=' in output, got:\n%s", out)
	}
}

func TestExport_SortKeys(t *testing.T) {
	entries := []Entry{
		{Key: "ZEBRA", Value: "z"},
		{Key: "ALPHA", Value: "a"},
		{Key: "MIDDLE", Value: "m"},
	}
	out, err := Export(entries, ExportOptions{Format: FormatDotEnv, SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	idxAlpha := strings.Index(out, "ALPHA")
	idxMiddle := strings.Index(out, "MIDDLE")
	idxZebra := strings.Index(out, "ZEBRA")
	if !(idxAlpha < idxMiddle && idxMiddle < idxZebra) {
		t.Errorf("keys not sorted: positions alpha=%d middle=%d zebra=%d", idxAlpha, idxMiddle, idxZebra)
	}
}

func TestExport_MaskSecrets(t *testing.T) {
	entries := []Entry{
		{Key: "API_SECRET", Value: "supersecret"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	out, err := Export(entries, ExportOptions{
		Format:      FormatDotEnv,
		MaskSecrets: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "supersecret") {
		t.Errorf("expected secret to be masked, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected non-secret value to remain, got:\n%s", out)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: "bar"}}
	_, err := Export(entries, ExportOptions{Format: "xml"})
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
