package envfile

import (
	"strings"
	"testing"
)

func TestGenerateTemplate_NonSecret(t *testing.T) {
	entries := []Entry{{Key: "APP_ENV", Value: "production"}}
	lines := GenerateTemplate(entries, nil)
	found := false
	for _, l := range lines {
		if strings.HasPrefix(l, "APP_ENV=") {
			found = true
			if !strings.Contains(l, "production") {
				t.Error("expected value preserved for non-secret key")
			}
		}
	}
	if !found {
		t.Error("APP_ENV line not found in template")
	}
}

func TestGenerateTemplate_SecretKey(t *testing.T) {
	masker := NewMasker(nil)
	entries := []Entry{{Key: "SECRET_KEY", Value: "supersecret"}}
	lines := GenerateTemplate(entries, masker)
	for _, l := range lines {
		if strings.Contains(l, "supersecret") {
			t.Error("secret value should not appear in template")
		}
	}
	hasBlank := false
	for _, l := range lines {
		if l == "SECRET_KEY=" {
			hasBlank = true
		}
	}
	if !hasBlank {
		t.Error("expected blank placeholder for SECRET_KEY")
	}
}

func TestApplyTemplate_FillsValues(t *testing.T) {
	template := []Entry{{Key: "DB_HOST", Value: ""}}
	values := map[string]string{"DB_HOST": "localhost"}
	result, errs := ApplyTemplate(template, values)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if result[0].Value != "localhost" {
		t.Errorf("expected localhost, got %s", result[0].Value)
	}
}

func TestApplyTemplate_MissingRequired(t *testing.T) {
	template := []Entry{{Key: "API_SECRET", Value: ""}}
	_, errs := ApplyTemplate(template, map[string]string{})
	if len(errs) == 0 {
		t.Fatal("expected error for missing required key")
	}
	if !strings.Contains(errs[0].Error(), "API_SECRET") {
		t.Errorf("unexpected error message: %v", errs[0])
	}
}

func TestApplyTemplate_UsesDefault(t *testing.T) {
	template := []Entry{{Key: "LOG_LEVEL", Value: "info"}}
	result, errs := ApplyTemplate(template, map[string]string{})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if result[0].Value != "info" {
		t.Errorf("expected default info, got %s", result[0].Value)
	}
}

func TestTemplateKeys_ReturnsEmpty(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: ""},
		{Key: "B", Value: "filled"},
		{Key: "C", Value: ""},
	}
	keys := TemplateKeys(entries)
	if len(keys) != 2 {
		t.Errorf("expected 2 empty keys, got %d", len(keys))
	}
}
