package envfile

import (
	"testing"
)

func TestValidate_ValidEntries(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "PORT", Value: "8080"},
		{Key: "DB_URL", Value: "postgres://localhost/db"},
	}
	result := Validate(entries)
	if !result.Valid() {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
}

func TestValidate_DuplicateKeys(t *testing.T) {
	entries := []Entry{
		{Key: "PORT", Value: "8080"},
		{Key: "PORT", Value: "9090"},
	}
	result := Validate(entries)
	if result.Valid() {
		t.Fatal("expected validation errors for duplicate keys")
	}
	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Key != "PORT" {
		t.Errorf("expected error on key PORT, got %q", result.Errors[0].Key)
	}
}

func TestValidate_InvalidKeyChars(t *testing.T) {
	entries := []Entry{
		{Key: "APP-NAME", Value: "bad"},
		{Key: "MY KEY", Value: "bad"},
	}
	result := Validate(entries)
	if result.Valid() {
		t.Fatal("expected validation errors for invalid key characters")
	}
	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(result.Errors))
	}
}

func TestValidate_KeyStartsWithDigit(t *testing.T) {
	entries := []Entry{
		{Key: "1BAD_KEY", Value: "value"},
	}
	result := Validate(entries)
	if result.Valid() {
		t.Fatal("expected error for key starting with digit")
	}
}

func TestValidate_EmptyKey(t *testing.T) {
	entries := []Entry{
		{Key: "", Value: "orphan"},
	}
	result := Validate(entries)
	if result.Valid() {
		t.Fatal("expected error for empty key")
	}
	if result.Errors[0].Message != "empty key" {
		t.Errorf("unexpected message: %s", result.Errors[0].Message)
	}
}

func TestValidateRequiredKeys_AllPresent(t *testing.T) {
	entries := []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "SECRET_KEY", Value: "abc"},
	}
	result := ValidateRequiredKeys(entries, []string{"APP_ENV", "SECRET_KEY"})
	if !result.Valid() {
		t.Errorf("expected valid, got: %v", result.Errors)
	}
}

func TestValidateRequiredKeys_Missing(t *testing.T) {
	entries := []Entry{
		{Key: "APP_ENV", Value: "production"},
	}
	result := ValidateRequiredKeys(entries, []string{"APP_ENV", "DATABASE_URL"})
	if result.Valid() {
		t.Fatal("expected missing key error")
	}
	if len(result.Errors) != 1 || result.Errors[0].Key != "DATABASE_URL" {
		t.Errorf("unexpected errors: %v", result.Errors)
	}
}
