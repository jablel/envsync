package envfile

import (
	"testing"
)

var resolveBase = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "DB_PORT", Value: "5432"},
}

var resolveOverride = []Entry{
	{Key: "DB_HOST", Value: "prod-db.example.com"},
	{Key: "NEW_KEY", Value: "new_value"},
}

func TestResolve_PreferOverride(t *testing.T) {
	opts := DefaultResolveOptions()
	r, err := Resolve(resolveBase, resolveOverride, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := toMap(r.Entries)
	if m["DB_HOST"] != "prod-db.example.com" {
		t.Errorf("expected override value, got %q", m["DB_HOST"])
	}
	if m["NEW_KEY"] != "new_value" {
		t.Errorf("expected new key, got %q", m["NEW_KEY"])
	}
	if r.Applied != 2 {
		t.Errorf("expected 2 applied, got %d", r.Applied)
	}
}

func TestResolve_PreferBase(t *testing.T) {
	opts := DefaultResolveOptions()
	opts.Strategy = ResolvePreferBase
	r, err := Resolve(resolveBase, resolveOverride, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := toMap(r.Entries)
	if m["DB_HOST"] != "localhost" {
		t.Errorf("expected base value, got %q", m["DB_HOST"])
	}
	if len(r.Conflicts) != 1 || r.Conflicts[0] != "DB_HOST" {
		t.Errorf("expected conflict on DB_HOST, got %v", r.Conflicts)
	}
}

func TestResolve_ErrorStrategy(t *testing.T) {
	opts := DefaultResolveOptions()
	opts.Strategy = ResolveError
	_, err := Resolve(resolveBase, resolveOverride, opts)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestResolve_EnvSuffix(t *testing.T) {
	override := []Entry{
		{Key: "DB_HOST_PROD", Value: "prod-db.example.com"},
		{Key: "APP_NAME_PROD", Value: "myapp-prod"},
	}
	opts := DefaultResolveOptions()
	opts.EnvSuffix = "_PROD"
	opts.StripSuffix = true
	r, err := Resolve(resolveBase, override, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := toMap(r.Entries)
	if m["DB_HOST"] != "prod-db.example.com" {
		t.Errorf("expected prod override, got %q", m["DB_HOST"])
	}
	if m["APP_NAME"] != "myapp-prod" {
		t.Errorf("expected prod app name, got %q", m["APP_NAME"])
	}
}

func TestResolve_NoOverride(t *testing.T) {
	opts := DefaultResolveOptions()
	r, err := Resolve(resolveBase, nil, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Entries) != len(resolveBase) {
		t.Errorf("expected %d entries, got %d", len(resolveBase), len(r.Entries))
	}
	if r.Applied != 0 {
		t.Errorf("expected 0 applied, got %d", r.Applied)
	}
}

func toMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
