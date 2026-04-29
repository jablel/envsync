package envfile

import (
	"errors"
	"strings"
	"testing"
)

func TestTransform_UppercaseKeys(t *testing.T) {
	entries := []Entry{
		{Key: "db_host", Value: "localhost"},
		{Key: "app_port", Value: "8080"},
	}
	out, err := Transform(entries, TransformOptions{}, UppercaseKeys())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "DB_HOST" || out[1].Key != "APP_PORT" {
		t.Errorf("keys not uppercased: %v", out)
	}
}

func TestTransform_PrefixKeys(t *testing.T) {
	entries := []Entry{{Key: "HOST", Value: "localhost"}}
	out, err := Transform(entries, TransformOptions{}, PrefixKeys("APP_"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %s", out[0].Key)
	}
}

func TestTransform_TrimValues(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: "  bar  "}}
	out, err := Transform(entries, TransformOptions{}, TrimValues())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "bar" {
		t.Errorf("expected 'bar', got %q", out[0].Value)
	}
}

func TestTransform_ChainedFuncs(t *testing.T) {
	entries := []Entry{{Key: "db_pass", Value: "  secret  "}}
	out, err := Transform(entries, TransformOptions{}, UppercaseKeys(), TrimValues())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "DB_PASS" {
		t.Errorf("key not uppercased")
	}
	if out[0].Value != "secret" {
		t.Errorf("value not trimmed")
	}
}

func TestTransform_ErrorAborts(t *testing.T) {
	entries := []Entry{{Key: "HOST", Value: "v"}}
	badFn := func(e Entry) (Entry, error) { return e, errors.New("boom") }
	_, err := Transform(entries, TransformOptions{}, badFn)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "HOST") {
		t.Errorf("error should mention key, got: %v", err)
	}
}

func TestTransform_SkipErrors(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	badFn := func(e Entry) (Entry, error) {
		if e.Key == "A" {
			return e, errors.New("skip me")
		}
		e.Value = "ok"
		return e, nil
	}
	out, err := Transform(entries, TransformOptions{SkipErrors: true}, badFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// A should be unchanged, B should be transformed
	if out[0].Value != "1" {
		t.Errorf("A should be unchanged, got %q", out[0].Value)
	}
	if out[1].Value != "ok" {
		t.Errorf("B should be transformed, got %q", out[1].Value)
	}
}

func TestTransform_MaskTransform(t *testing.T) {
	m := NewMasker(nil)
	entries := []Entry{
		{Key: "API_SECRET", Value: "topsecret"},
		{Key: "HOST", Value: "localhost"},
	}
	out, err := Transform(entries, TransformOptions{}, MaskTransform(m))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value == "topsecret" {
		t.Error("sensitive value should have been masked")
	}
	if out[1].Value != "localhost" {
		t.Error("non-sensitive value should be unchanged")
	}
}

func TestTransform_EmptyFuncs(t *testing.T) {
	entries := []Entry{{Key: "K", Value: "v"}}
	out, err := Transform(entries, TransformOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Key != "K" {
		t.Error("entries should be returned unchanged")
	}
}
