package envfile

import (
	"testing"
)

func TestClone_CopiesNewKeys(t *testing.T) {
	src := []Entry{{Key: "FOO", Value: "1"}, {Key: "BAR", Value: "2"}}
	dst := []Entry{}

	out, res, err := Clone(src, dst, CloneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
}

func TestClone_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := []Entry{{Key: "FOO", Value: "new"}}
	dst := []Entry{{Key: "FOO", Value: "old"}}

	out, res, err := Clone(src, dst, CloneOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "old" {
		t.Errorf("expected old value to be preserved, got %q", out[0].Value)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "FOO" {
		t.Errorf("expected FOO in skipped, got %v", res.Skipped)
	}
}

func TestClone_OverwriteExistingKeys(t *testing.T) {
	src := []Entry{{Key: "FOO", Value: "new"}}
	dst := []Entry{{Key: "FOO", Value: "old"}}

	out, res, err := Clone(src, dst, CloneOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "new" {
		t.Errorf("expected new value, got %q", out[0].Value)
	}
	if len(res.Copied) != 1 {
		t.Errorf("expected 1 copied, got %d", len(res.Copied))
	}
}

func TestClone_KeyFilter(t *testing.T) {
	src := []Entry{{Key: "FOO", Value: "1"}, {Key: "BAR", Value: "2"}, {Key: "BAZ", Value: "3"}}
	dst := []Entry{}

	opts := CloneOptions{
		KeyFilter: func(k string) bool { return k == "FOO" || k == "BAZ" },
	}
	out, res, err := Clone(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %v", res.Copied)
	}
}

func TestClone_StripSecrets(t *testing.T) {
	src := []Entry{
		{Key: "APP_SECRET_KEY", Value: "supersecret"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	dst := []Entry{}

	out, res, err := Clone(src, dst, CloneOptions{StripSecrets: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Key == "APP_SECRET_KEY" && e.Value != "" {
			t.Errorf("expected secret to be stripped, got %q", e.Value)
		}
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("expected non-secret value preserved, got %q", e.Value)
		}
	}
	if len(res.Stripped) == 0 {
		t.Errorf("expected at least one stripped key")
	}
}

func TestClone_NilSrcReturnsError(t *testing.T) {
	_, _, err := Clone(nil, []Entry{}, CloneOptions{})
	if err == nil {
		t.Error("expected error for nil src, got nil")
	}
}
