package envfile

import (
	"testing"
)

func schemaEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "PORT", Value: "8080"},
		{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
	}
}

func TestValidateSchema_AllValid(t *testing.T) {
	fields := []SchemaField{
		{Key: "APP_NAME", Required: true},
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
		{Key: "DATABASE_URL", Required: true},
	}
	res, err := ValidateSchema(schemaEntries(), fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HasErrors() {
		t.Errorf("expected no errors, got: %s", res.Summary())
	}
	if len(res.Valid) != 3 {
		t.Errorf("expected 3 valid keys, got %d", len(res.Valid))
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	fields := []SchemaField{
		{Key: "APP_NAME", Required: true},
		{Key: "SECRET_KEY", Required: true},
	}
	res, err := ValidateSchema(schemaEntries(), fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY missing, got %v", res.Missing)
	}
}

func TestValidateSchema_InvalidPattern(t *testing.T) {
	fields := []SchemaField{
		{Key: "PORT", Required: true, Pattern: `^[a-z]+$`},
	}
	res, err := ValidateSchema(schemaEntries(), fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Invalid) != 1 || res.Invalid[0] != "PORT" {
		t.Errorf("expected PORT invalid, got %v", res.Invalid)
	}
}

func TestValidateSchema_ExtraKeys(t *testing.T) {
	fields := []SchemaField{
		{Key: "APP_NAME", Required: true},
	}
	res, err := ValidateSchema(schemaEntries(), fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Extra) != 2 {
		t.Errorf("expected 2 extra keys, got %d: %v", len(res.Extra), res.Extra)
	}
}

func TestValidateSchema_BadPatternReturnsError(t *testing.T) {
	fields := []SchemaField{
		{Key: "APP_NAME", Required: true, Pattern: `[invalid(`},
	}
	_, err := ValidateSchema(schemaEntries(), fields)
	if err == nil {
		t.Error("expected error for invalid regex pattern, got nil")
	}
}

func TestSchemaResult_Summary(t *testing.T) {
	res := SchemaResult{
		Missing: []string{"SECRET_KEY"},
		Invalid: []string{"PORT"},
		Extra:   []string{"UNUSED"},
		Valid:   []string{"APP_NAME"},
	}
	summary := res.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	if !res.HasErrors() {
		t.Error("expected HasErrors to be true")
	}
}
