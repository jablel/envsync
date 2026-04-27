package envfile

import (
	"os"
	"testing"
)

func TestInterpolate_BasicSubstitution(t *testing.T) {
	entries := []Entry{
		{Key: "BASE_URL", Value: "https://example.com"},
		{Key: "API_URL", Value: "${BASE_URL}/api"},
	}

	result, err := Interpolate(entries, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[1].Value != "https://example.com/api" {
		t.Errorf("expected expanded URL, got %q", result[1].Value)
	}
}

func TestInterpolate_DollarSignStyle(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "DSN", Value: "postgres://$HOST/db"},
	}

	result, err := Interpolate(entries, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[1].Value != "postgres://localhost/db" {
		t.Errorf("got %q", result[1].Value)
	}
}

func TestInterpolate_FallbackToOSEnv(t *testing.T) {
	os.Setenv("OS_VAR", "from-os")
	defer os.Unsetenv("OS_VAR")

	entries := []Entry{
		{Key: "VALUE", Value: "${OS_VAR}-suffix"},
	}

	opts := DefaultInterpolateOptions()
	opts.AllowEnvFallback = true

	result, err := Interpolate(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Value != "from-os-suffix" {
		t.Errorf("expected OS fallback, got %q", result[0].Value)
	}
}

func TestInterpolate_NoFallback_LeavesUnresolved(t *testing.T) {
	entries := []Entry{
		{Key: "VALUE", Value: "${UNDEFINED_XYZ}"},
	}

	opts := DefaultInterpolateOptions()
	opts.AllowEnvFallback = false
	opts.FailOnMissing = false

	result, err := Interpolate(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Unresolved reference should remain as-is.
	if result[0].Value != "${UNDEFINED_XYZ}" {
		t.Errorf("expected unchanged value, got %q", result[0].Value)
	}
}

func TestInterpolate_FailOnMissing(t *testing.T) {
	entries := []Entry{
		{Key: "VALUE", Value: "${DEFINITELY_MISSING}"},
	}

	opts := DefaultInterpolateOptions()
	opts.AllowEnvFallback = false
	opts.FailOnMissing = true

	_, err := Interpolate(entries, opts)
	if err == nil {
		t.Fatal("expected error for missing variable, got nil")
	}
}

func TestInterpolate_PreservesComments(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "db", Comment: "# database host"},
		{Key: "URL", Value: "${HOST}:5432", Comment: "# connection url"},
	}

	result, err := Interpolate(entries, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[1].Comment != "# connection url" {
		t.Errorf("comment not preserved: %q", result[1].Comment)
	}
	if result[1].Value != "db:5432" {
		t.Errorf("expected db:5432, got %q", result[1].Value)
	}
}
